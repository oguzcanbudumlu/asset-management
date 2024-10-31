package e2e

import (
	"asset-management/services/asset-api/deposit"
	"asset-management/services/asset-api/dto"
	"asset-management/services/asset-api/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDepositEndpoint(t *testing.T) {
	db, cleanup := util.SetupTestContainer(t)
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
