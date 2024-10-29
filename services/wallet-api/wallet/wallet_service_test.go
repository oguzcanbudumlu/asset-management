package wallet

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

// MockWalletRepository is a mock implementation of the WalletRepository interface
type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) CreateWallet(wallet *Wallet) error {
	args := m.Called(wallet)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWallets() ([]Wallet, error) {
	args := m.Called()
	return args.Get(0).([]Wallet), args.Error(1)
}

func (m *MockWalletRepository) DeleteWallet(network, address string) error {
	args := m.Called(network, address)
	return args.Error(0)
}

func (m *MockWalletRepository) GetWallet(network, address string) (*Wallet, error) {
	args := m.Called(network, address)
	if wallet, ok := args.Get(0).(*Wallet); ok {
		return wallet, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestServiceCreateWallet(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	wallet := &Wallet{ /* initialize your wallet */ }

	// Set up expectation
	mockRepo.On("CreateWallet", wallet).Return(nil)

	// Call the method
	err := service.CreateWallet(wallet)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestServiceGetWallet(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	network := "test-network"
	address := "test-address"
	wallet := &Wallet{ /* initialize your wallet */ }

	// Set up expectation
	mockRepo.On("GetWallet", network, address).Return(wallet, nil)

	// Call the method
	result, err := service.GetWallet(network, address)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, wallet, result)
	mockRepo.AssertExpectations(t)
}

func TestServiceDeleteWallet(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	network := "test-network"
	address := "test-address"

	// Set up expectation
	mockRepo.On("DeleteWallet", network, address).Return(nil)

	// Call the method
	err := service.DeleteWallet(network, address)

	// Assert
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestCreateWallet_Failure(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	wallet := &Wallet{ /* initialize your wallet */ }

	// Set up expectation with an error
	mockRepo.On("CreateWallet", wallet).Return(errors.New("database error"))

	// Call the method
	err := service.CreateWallet(wallet)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "database error", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestGetWallet_Failure(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	network := "test-network"
	address := "test-address"

	// Set up expectation with an error
	mockRepo.On("GetWallet", network, address).Return(nil, errors.New("wallet not found"))

	// Call the method
	result, err := service.GetWallet(network, address)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "wallet not found", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestDeleteWallet_Failure(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	service := NewWalletService(mockRepo)

	network := "test-network"
	address := "test-address"

	// Set up expectation with an error
	mockRepo.On("DeleteWallet", network, address).Return(errors.New("delete failed"))

	// Call the method
	err := service.DeleteWallet(network, address)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "delete failed", err.Error())
	mockRepo.AssertExpectations(t)
}
