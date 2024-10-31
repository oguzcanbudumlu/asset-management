package scheduled_process_test

import (
	"asset-management/internal/schedule/scheduled_process"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockProcessService struct {
	mock.Mock
}

func (m *MockProcessService) Process(scheduledTransactionID int) error {
	args := m.Called(scheduledTransactionID)
	return args.Error(0)
}

func TestProcessController_Process_Success(t *testing.T) {
	mockService := new(MockProcessService)
	controller := scheduled_process.NewProcessController(mockService)

	app := fiber.New()
	app.Post("/scheduled-transaction/:id/process", controller.Process)

	// Mock successful service response
	mockService.On("Process", 123).Return(nil)

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction/123/process", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Transaction processed successfully", response["message"])

	mockService.AssertExpectations(t)
}

func TestProcessController_Process_InvalidID(t *testing.T) {
	mockService := new(MockProcessService)
	controller := scheduled_process.NewProcessController(mockService)

	app := fiber.New()
	app.Post("/scheduled-transaction/:id/process", controller.Process)

	// Create HTTP request with invalid transaction ID
	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction/invalid/process", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "Invalid transaction ID", response["error"])
}

func TestProcessController_Process_Failure(t *testing.T) {
	mockService := new(MockProcessService)
	controller := scheduled_process.NewProcessController(mockService)

	app := fiber.New()
	app.Post("/scheduled-transaction/:id/process", controller.Process)

	// Mock error in service response
	mockService.On("Process", 123).Return(errors.New("service error"))

	// Create HTTP request
	req := httptest.NewRequest(http.MethodPost, "/scheduled-transaction/123/process", nil)
	resp, err := app.Test(req)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]string
	_ = json.NewDecoder(resp.Body).Decode(&response)
	assert.Equal(t, "failed to process transaction: service error", response["error"])

	mockService.AssertExpectations(t)
}
