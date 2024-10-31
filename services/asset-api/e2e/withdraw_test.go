package e2e_test

import (
	"asset-management/services/asset-api/util"
	"asset-management/services/asset-api/withdraw"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockValidationAdapter for simulating wallet validation behavior
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

func setup(t *testing.T) (*fiber.App, *sql.DB, func(), *MockValidationAdapter) {
	// Setup PostgreSQL test container and initialize schema
	db, cleanup := util.SetupTestContainer(t)

	// Create the mock validation adapter
	mockValidation := new(MockValidationAdapter)

	// Set up repository, service, and controller with mockValidation
	repo := withdraw.NewRepository(db)
	service := withdraw.NewService(repo, mockValidation)
	controller := withdraw.NewController(service)

	// Setup Fiber app with the withdraw route
	app := fiber.New()
	app.Post("/withdraw", controller.Withdraw)

	return app, db, cleanup, mockValidation
}

func TestE2E_Withdraw_Success_WithMock(t *testing.T) {
	app, db, cleanup, mockValidation := setup(t)
	defer cleanup()

	err := util.InsertBalance(db, "0x123abc456def", "Ethereum", 500.00)
	assert.NoError(t, err)
	// Mock successful wallet validation
	mockValidation.On("One", "0x123abc456def", "Ethereum").Return(nil)

	// Define withdraw request payload
	reqPayload := withdraw.Request{
		WalletAddress: "0x123abc456def",
		Network:       "Ethereum",
		Amount:        100.50,
	}
	reqBody, _ := json.Marshal(reqPayload)

	// Create HTTP request to test the /withdraw endpoint
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Verify the updated balance
	var newBalance float64
	err = db.QueryRow(`
		SELECT balance FROM balance 
		WHERE wallet_address = $1 AND network = $2`, "0x123abc456def", "Ethereum").Scan(&newBalance)
	assert.NoError(t, err)
	assert.Equal(t, 399.50, newBalance)

	// Assert that the mock was called as expected
	mockValidation.AssertExpectations(t)
}

func TestE2E_Withdraw_InsufficientBalance_WithMock(t *testing.T) {
	app, db, cleanup, mockValidation := setup(t)
	defer cleanup()

	_ = util.InsertBalance(db, "0x123abc456def", "Ethereum", 50.00)

	// Mock successful wallet validation
	mockValidation.On("One", "0x123abc456def", "Ethereum").Return(nil)

	// Define withdraw request payload with an amount greater than the balance
	reqPayload := withdraw.Request{
		WalletAddress: "0x123abc456def",
		Network:       "Ethereum",
		Amount:        100.50,
	}
	reqBody, _ := json.Marshal(reqPayload)

	// Create HTTP request to test the /withdraw endpoint
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Check response message for "insufficient balance"
	var errorResponse struct {
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.Equal(t, "withdraw transaction failed: insufficient balance", errorResponse.Message)

	// Assert that the mock was called as expected
	mockValidation.AssertExpectations(t)
}

func TestE2E_Withdraw_WalletNotFound_WithMock(t *testing.T) {
	app, _, cleanup, mockValidation := setup(t)
	defer cleanup()

	// Mock wallet validation to simulate wallet not found
	mockValidation.On("One", "0xUnknownWallet", "Ethereum").Return(errors.New("wallet not found"))

	// Define withdraw request payload for a wallet that doesn’t exist in the database
	reqPayload := withdraw.Request{
		WalletAddress: "0xUnknownWallet",
		Network:       "Ethereum",
		Amount:        100.50,
	}
	reqBody, _ := json.Marshal(reqPayload)

	// Create HTTP request to test the /withdraw endpoint
	req := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	// Check response message for "wallet not found"
	var errorResponse struct {
		Message string `json:"message"`
	}
	_ = json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.Equal(t, "wallet validation failed: wallet not found", errorResponse.Message)

	// Assert that the mock was called as expected
	mockValidation.AssertExpectations(t)
}
