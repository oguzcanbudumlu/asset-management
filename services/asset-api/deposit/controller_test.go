// controller_test.go

package deposit_test

import (
	"asset-management/services/asset-api/deposit"
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Mock service
type mockDepositService struct{ mock.Mock }

func (m *mockDepositService) Deposit(walletAddress, network string, amount float64) (float64, error) {
	args := m.Called(walletAddress, network, amount)
	return args.Get(0).(float64), args.Error(1)
}

func TestDepositController_Success(t *testing.T) {
	service := new(mockDepositService)
	service.On("Deposit", "0x123abc456def", "Ethereum", 100.50).Return(1500.75, nil)

	controller := deposit.NewController(service)

	app := fiber.New()
	app.Post("/deposit", controller.Deposit)

	req := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader(`{"wallet_address":"0x123abc456def","network":"Ethereum","amount":100.50}`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}

func TestDepositController_InvalidPayload(t *testing.T) {
	service := new(mockDepositService)
	controller := deposit.NewController(service)

	app := fiber.New()
	app.Post("/deposit", controller.Deposit)

	req := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader(`invalid json`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDepositController_ServiceError(t *testing.T) {
	service := new(mockDepositService)
	service.On("Deposit", "0x123abc456def", "Ethereum", 100.50).Return(0.0, errors.New("deposit error"))

	controller := deposit.NewController(service)

	app := fiber.New()
	app.Post("/deposit", controller.Deposit)

	req := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader(`{"wallet_address":"0x123abc456def","network":"Ethereum","amount":100.50}`))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req, -1)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
