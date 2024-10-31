package withdraw_test

import (
	"asset-management/services/asset-api/common/dto"
	"asset-management/services/asset-api/withdraw"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

type MockService struct {
	mock.Mock
}

func (m *MockService) Withdraw(walletAddress, network string, amount float64) error {
	args := m.Called(walletAddress, network, amount)
	return args.Error(0)
}

func TestController_Withdraw_Success(t *testing.T) {
	app := fiber.New()
	mockService := new(MockService)
	controller := withdraw.NewController(mockService)

	app.Post("/withdraw", controller.Withdraw)

	// Arrange
	req := withdraw.Request{
		WalletAddress: "0x123abc456def",
		Network:       "Ethereum",
		Amount:        100.50,
	}
	mockService.On("Withdraw", req.WalletAddress, req.Network, req.Amount).Return(nil)

	body, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, _ := app.Test(request)

	// Assert
	assert.Equal(t, http.StatusOK, response.StatusCode)
	mockService.AssertExpectations(t)
}

func TestController_Withdraw_InvalidPayload(t *testing.T) {
	app := fiber.New()
	mockService := new(MockService)
	controller := withdraw.NewController(mockService)

	app.Post("/withdraw", controller.Withdraw)

	// Act
	request := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer([]byte(`invalid json`)))
	request.Header.Set("Content-Type", "application/json")
	response, _ := app.Test(request)

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	var errorResponse dto.ErrorResponse
	_ = json.NewDecoder(response.Body).Decode(&errorResponse)
	assert.Equal(t, "Invalid request payload", errorResponse.Message)
}

func TestController_Withdraw_ServiceError(t *testing.T) {
	app := fiber.New()
	mockService := new(MockService)
	controller := withdraw.NewController(mockService)

	app.Post("/withdraw", controller.Withdraw)

	// Arrange
	req := withdraw.Request{
		WalletAddress: "0x123abc456def",
		Network:       "Ethereum",
		Amount:        100.50,
	}
	mockService.On("Withdraw", req.WalletAddress, req.Network, req.Amount).Return(errors.New("insufficient balance"))

	body, _ := json.Marshal(req)
	request := httptest.NewRequest(http.MethodPost, "/withdraw", bytes.NewBuffer(body))
	request.Header.Set("Content-Type", "application/json")
	response, _ := app.Test(request)

	// Assert
	assert.Equal(t, fiber.StatusBadRequest, response.StatusCode)

	var errorResponse dto.ErrorResponse
	_ = json.NewDecoder(response.Body).Decode(&errorResponse)
	assert.Equal(t, "insufficient balance", errorResponse.Message)
}
