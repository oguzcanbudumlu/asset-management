package e2e_test

import (
	"asset-management/services/asset-api/transaction"
	"asset-management/services/asset-api/util"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Setup function to initialize Fiber app with repository, service, and controller
func setupE2EApp(t *testing.T) (*fiber.App, func(), *MockValidationAdapter) {
	// Setup PostgreSQL test container and initialize schema
	db, cleanup := util.SetupTestContainer(t)

	// Create repository, mock validation adapter, and service
	repo := transaction.NewCreateRepository(db)
	mockValidator := new(MockValidationAdapter)
	service := transaction.NewCreateService(repo, mockValidator)
	controller := transaction.NewCreateController(service)

	// Setup Fiber app with the transaction route
	app := fiber.New()
	app.Post("/scheduled-transaction", controller.Create)

	return app, cleanup, mockValidator
}

func TestE2E_CreateTransaction_Success(t *testing.T) {
	app, cleanup, mockValidator := setupE2EApp(t)
	defer cleanup()

	// Mock successful wallet validation
	mockValidator.On("Both", "wallet123", "wallet456", "mainnet").Return(nil)

	// Define valid transaction request payload
	reqPayload := transaction.Request{
		From:          "wallet123",
		To:            "wallet456",
		Network:       "mainnet",
		Amount:        100.50,
		ScheduledTime: "2023-12-31T12:00:00Z",
	}
	reqBody, _ := json.Marshal(reqPayload)

	// Create HTTP request to test the /scheduled-transaction endpoint
	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request to Fiber app and check response
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	// Parse and verify the response
	var response map[string]int
	json.NewDecoder(resp.Body).Decode(&response)
	assert.NotZero(t, response["transaction_id"])

	// Assert that the mock validation was called as expected
	mockValidator.AssertExpectations(t)
}

func TestE2E_CreateTransaction_InvalidAmount(t *testing.T) {
	app, cleanup, _ := setupE2EApp(t)
	defer cleanup()

	// Define transaction request payload with an invalid amount
	reqPayload := transaction.Request{
		From:          "wallet123",
		To:            "wallet456",
		Network:       "mainnet",
		Amount:        0, // Invalid amount
		ScheduledTime: "2023-12-31T12:00:00Z",
	}
	reqBody, _ := json.Marshal(reqPayload)

	// Create HTTP request to test the /scheduled-transaction endpoint
	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request and check response for validation error
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var errorResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.Equal(t, "amount must be greater than zero", errorResponse["error"])
}

func TestE2E_CreateTransaction_InvalidScheduledTimeFormat(t *testing.T) {
	app, cleanup, _ := setupE2EApp(t)
	defer cleanup()

	// Define transaction request payload with an invalid scheduled time format
	reqPayload := transaction.Request{
		From:          "wallet123",
		To:            "wallet456",
		Network:       "mainnet",
		Amount:        100.50,
		ScheduledTime: "invalid-time-format", // Invalid time format
	}
	reqBody, _ := json.Marshal(reqPayload)

	// Create HTTP request to test the /scheduled-transaction endpoint
	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request and check response for invalid time format error
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var errorResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.Equal(t, "Invalid scheduled time format", errorResponse["error"])
}

func TestE2E_CreateTransaction_ValidationFailed(t *testing.T) {
	app, cleanup, mockValidator := setupE2EApp(t)
	defer cleanup()

	// Mock validation failure for wallets
	mockValidator.On("Both", "wallet123", "wallet456", "mainnet").Return(errors.New("wallet validation failed"))

	// Define transaction request payload
	reqPayload := transaction.Request{
		From:          "wallet123",
		To:            "wallet456",
		Network:       "mainnet",
		Amount:        100.50,
		ScheduledTime: "2023-12-31T12:00:00Z",
	}
	reqBody, _ := json.Marshal(reqPayload)

	// Create HTTP request to test the /scheduled-transaction endpoint
	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	// Send request and check response for validation error
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var errorResponse map[string]string
	json.NewDecoder(resp.Body).Decode(&errorResponse)
	assert.Equal(t, "wallet validation failed", errorResponse["error"])

	// Assert that the mock validation was called as expected
	mockValidator.AssertExpectations(t)
}
