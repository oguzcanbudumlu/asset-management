package scheduled_test

import (
	"asset-management/internal/schedule"
	"asset-management/services/asset-api/scheduled"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type mockNextService struct {
	mock.Mock
}

func (m *mockNextService) GetNextMinuteTransactions() ([]schedule.ScheduledTransaction, error) {
	args := m.Called()
	if transactions, ok := args.Get(0).([]schedule.ScheduledTransaction); ok {
		return transactions, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestNextController_GetNextMinuteTransactions_Success(t *testing.T) {
	mockService := new(mockNextService)
	controller := scheduled.NewNextController(mockService)

	app := fiber.New()
	app.Get("/scheduled-transaction/next", controller.GetNextMinuteTransactions)

	// Mock successful response
	transactions := []schedule.ScheduledTransaction{
		{ID: 1, FromWallet: "wallet123", ToWallet: "wallet456", Network: "mainnet", Amount: 100.50},
	}
	mockService.On("GetNextMinuteTransactions").Return(transactions, nil)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/scheduled-transaction/next", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response []schedule.ScheduledTransaction
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, transactions, response)

	mockService.AssertExpectations(t)
}

func TestNextController_GetNextMinuteTransactions_Error(t *testing.T) {
	mockService := new(mockNextService)
	controller := scheduled.NewNextController(mockService)

	app := fiber.New()
	app.Get("/scheduled-transaction/next", controller.GetNextMinuteTransactions)

	// Mock error response
	mockService.On("GetNextMinuteTransactions").Return(nil, errors.New("service error"))

	// Create HTTP request
	req := httptest.NewRequest(http.MethodGet, "/scheduled-transaction/next", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Failed to retrieve transactions", response["error"])

	mockService.AssertExpectations(t)
}
