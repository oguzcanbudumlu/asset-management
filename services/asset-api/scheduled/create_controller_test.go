package scheduled_test

import (
	"asset-management/services/asset-api/scheduled"
	"bytes"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type MockCreateService struct {
	mock.Mock
}

func (m *MockCreateService) Create(fromWallet, toWallet, network string, amount float64, scheduledTime time.Time) (int, error) {
	args := m.Called(fromWallet, toWallet, network, amount, scheduledTime)
	return args.Int(0), args.Error(1)
}

func TestCreateController_Success(t *testing.T) {
	mockService := new(MockCreateService)
	controller := scheduled.NewCreateController(mockService)
	app := fiber.New()
	app.Post("/scheduled-transaction", controller.Create)

	reqPayload := scheduled.Request{
		From:          "wallet123",
		To:            "wallet456",
		Network:       "mainnet",
		Amount:        100.50,
		ScheduledTime: "2023-12-31T12:00:00Z",
	}
	reqBody, _ := json.Marshal(reqPayload)

	mockService.On("Create", "wallet123", "wallet456", "mainnet", 100.50, time.Date(2023, 12, 31, 12, 0, 0, 0, time.UTC)).Return(123, nil)

	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	var response map[string]int
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, 123, response["transaction_id"])

	mockService.AssertExpectations(t)
}

func TestCreateController_InvalidScheduledTime(t *testing.T) {
	mockService := new(MockCreateService)
	controller := scheduled.NewCreateController(mockService)
	app := fiber.New()
	app.Post("/scheduled-transaction", controller.Create)

	reqPayload := scheduled.Request{
		From:          "wallet123",
		To:            "wallet456",
		Network:       "mainnet",
		Amount:        100.50,
		ScheduledTime: "invalid-time-format",
	}
	reqBody, _ := json.Marshal(reqPayload)

	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Invalid scheduled time format", response["error"])
}
