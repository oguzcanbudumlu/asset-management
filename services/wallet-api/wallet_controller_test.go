package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWalletService is a mock implementation of WalletService
type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) CreateWallet(wallet *Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletService) GetWallet(network, address string) (*Wallet, error) {
	args := m.Called(network, address)

	if wallet, ok := args.Get(0).(*Wallet); ok {
		return wallet, args.Error(1)
	}
	return nil, args.Error(1) // Ensure to return nil if wallet is not found
}

func (m *MockWalletService) DeleteWallet(network, address string) error {
	args := m.Called(network, address)
	return args.Error(0)
}

func TestCreateWallet(t *testing.T) {
	app := fiber.New()
	mockService := new(MockWalletService)
	controller := NewWalletController(mockService)

	app.Post("/wallet", controller.CreateWallet)

	wallet := Wallet{Address: "test_address", Network: "test_network"}
	body, _ := json.Marshal(wallet)

	t.Run("success", func(t *testing.T) {
		mockService.On("CreateWallet", &wallet).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid request", func(t *testing.T) {
		reqInvalid := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewBuffer([]byte("invalid json")))
		reqInvalid.Header.Set("Content-Type", "application/json")
		respInvalid, _ := app.Test(reqInvalid)

		assert.Equal(t, http.StatusBadRequest, respInvalid.StatusCode)
	})
}

func TestGetWallet(t *testing.T) {
	app := fiber.New()
	mockService := new(MockWalletService)
	controller := NewWalletController(mockService)

	app.Get("/wallet/:network/:address", controller.GetWallet)

	// Test successful retrieval
	t.Run("success", func(t *testing.T) {
		mockWallet := &Wallet{Address: "test_address", Network: "test_network"}
		mockService.On("GetWallet", "test_network", "test_address").Return(mockWallet, nil)

		req := httptest.NewRequest(http.MethodGet, "/wallet/test_network/test_address", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var returnedWallet Wallet
		json.NewDecoder(resp.Body).Decode(&returnedWallet)

		assert.Equal(t, mockWallet, &returnedWallet)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("GetWallet", "test_network", "nonexistent_address").Return(nil, errors.New("error"))

		reqNotFound := httptest.NewRequest(http.MethodGet, "/wallet/test_network/nonexistent_address", nil)
		respNotFound, _ := app.Test(reqNotFound)

		assert.Equal(t, http.StatusNotFound, respNotFound.StatusCode)
	})
}

func TestDeleteWallet(t *testing.T) {
	app := fiber.New()
	mockService := new(MockWalletService)
	controller := NewWalletController(mockService)

	app.Delete("/wallet/:network/:address", controller.DeleteWallet)

	t.Run("successful", func(t *testing.T) {
		mockService.On("DeleteWallet", "test_network", "test_address").Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/wallet/test_network/test_address", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("not found", func(t *testing.T) {
		mockService.On("DeleteWallet", "test_network", "nonexistent_address").Return(errors.New("error"))

		reqNotFound := httptest.NewRequest(http.MethodDelete, "/wallet/test_network/nonexistent_address", nil)
		respNotFound, _ := app.Test(reqNotFound)

		assert.Equal(t, http.StatusNotFound, respNotFound.StatusCode)
	})
}
