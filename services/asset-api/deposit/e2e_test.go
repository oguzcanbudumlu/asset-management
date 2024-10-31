package deposit_test

import (
	sql2 "asset-management/internal/sql"
	"asset-management/services/asset-api/common/dto"
	"asset-management/services/asset-api/deposit"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockValidationAdapter struct {
	mock.Mock
}

func (m *MockValidationAdapter) One(walletAddress, network string) error {
	args := m.Called(walletAddress, network)
	return args.Error(0)
}

func (m *MockValidationAdapter) Both(from, to, network string) error {
	args := m.Called(from, to, network)
	return args.Error(0)
}

func SetupTestContainer(t *testing.T) (*sql.DB, func()) {
	t.Helper()
	ctx := context.Background()

	// Setting up PostgreSQL container
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "testuser",
			"POSTGRES_PASSWORD": "testpass",
			"POSTGRES_DB":       "testdb",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}
	pgContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}

	// Retrieve the host and port for PostgreSQL
	host, err := pgContainer.Host(ctx)
	if err != nil {
		t.Fatalf("Failed to get container host: %v", err)
	}
	port, err := pgContainer.MappedPort(ctx, "5432")
	if err != nil {
		t.Fatalf("Failed to get container port: %v", err)
	}

	// Database connection setup
	dsn := "postgres://testuser:testpass@" + host + ":" + port.Port() + "/testdb?sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		t.Fatalf("Failed to connect to the database: %v", err)
	}

	// Setup schema
	_, err = db.Exec(sql2.CreateBalanceTableSQL)
	assert.NoError(t, err)

	_, err = db.Exec(sql2.CreateBalanceTableSQL)
	assert.NoError(t, err)

	// Cleanup function to terminate the container
	cleanup := func() {
		db.Close()
		pgContainer.Terminate(ctx)
	}

	return db, cleanup
}

func TestDepositEndpoint(t *testing.T) {
	db, cleanup := SetupTestContainer(t)
	defer cleanup()

	mockValidation := new(MockValidationAdapter)

	repo := deposit.NewRepository(db)
	service := deposit.NewService(mockValidation, repo)
	controller := deposit.NewController(service)

	app := fiber.New()
	app.Post("/deposit", controller.Deposit)

	// Seed database with initial balance
	_, err := db.Exec("INSERT INTO balance (wallet_address, network, balance) VALUES ($1, $2, $3)", "0x123abc456def", "Ethereum", 1000.00)
	if err != nil {
		t.Fatalf("failed to seed database: %s", err)
	}

	sendRequest := func(t *testing.T, reqBody deposit.Request, expectedStatus int) *http.Response {
		reqBytes, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/deposit", bytes.NewReader(reqBytes))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to make request: %s", err)
		}
		if resp.StatusCode != expectedStatus {
			t.Errorf("expected status %d but got %d", expectedStatus, resp.StatusCode)
		}
		return resp
	}

	// 1. Successful Deposit with Valid Wallet
	t.Run("Successful Deposit", func(t *testing.T) {
		// Mock successful validation
		mockValidation.On("One", "0x123abc456def", "Ethereum").Return(nil)

		reqBody := deposit.Request{
			WalletAddress: "0x123abc456def",
			Network:       "Ethereum",
			Amount:        100.50,
		}
		resp := sendRequest(t, reqBody, http.StatusOK)

		var response deposit.Response
		err := json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("failed to decode response: %s", err)
		}

		expectedNewBalance := 1100.50
		if response.NewBalance != expectedNewBalance {
			t.Errorf("expected new balance %f but got %f", expectedNewBalance, response.NewBalance)
		}

		mockValidation.AssertCalled(t, "One", "0x123abc456def", "Ethereum")
	})

	// 2. Wallet Validation Failure
	t.Run("Wallet Validation Failure", func(t *testing.T) {
		// Mock validation failure
		mockValidation.On("One", "0xinvalidwallet", "Ethereum").Return(fmt.Errorf("wallet not found"))

		reqBody := deposit.Request{
			WalletAddress: "0xinvalidwallet",
			Network:       "Ethereum",
			Amount:        100.50,
		}
		resp := sendRequest(t, reqBody, http.StatusBadRequest)

		var response dto.ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("failed to decode error response: %s", err)
		}
		if response.Message != "wallet validation failed: wallet not found" {
			t.Errorf("expected error 'wallet validation failed: wallet not found', got %s", response.Message)
		}

		mockValidation.AssertCalled(t, "One", "0xinvalidwallet", "Ethereum")
	})

	// 3. Invalid Network
	t.Run("Invalid Network", func(t *testing.T) {
		mockValidation.On("One", "0x123abc456def", "").Return(nil)

		reqBody := deposit.Request{
			WalletAddress: "0x123abc456def",
			Network:       "",
			Amount:        100.50,
		}
		resp := sendRequest(t, reqBody, http.StatusBadRequest)

		var response dto.ErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			t.Fatalf("failed to decode error response: %s", err)
		}
		if response.Message != "invalid input parameters" {
			t.Errorf("expected error 'invalid input parameters', got %s", response.Message)
		}

		mockValidation.AssertNotCalled(t, "One", "0x123abc456def", "")
	})
}
