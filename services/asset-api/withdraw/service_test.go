package withdraw_test

import (
	"errors"
	"testing"

	"asset-management/services/asset-api/withdraw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock for ValidationAdapter
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

// Mock for Repository
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Withdraw(walletAddress, network string, amount float64) error {
	args := m.Called(walletAddress, network, amount)
	return args.Error(0)
}

func TestWithdrawService_Success(t *testing.T) {
	mockValidator := new(MockValidationAdapter)
	mockRepo := new(MockRepository)

	// Arrange
	service := withdraw.NewService(mockRepo, mockValidator)
	mockValidator.On("One", "0x123abc456def", "Ethereum").Return(nil)
	mockRepo.On("Withdraw", "0x123abc456def", "Ethereum", 100.50).Return(nil)

	// Act
	err := service.Withdraw("0x123abc456def", "Ethereum", 100.50)

	// Assert
	assert.NoError(t, err)
	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestWithdrawService_InvalidInput(t *testing.T) {
	mockValidator := new(MockValidationAdapter)
	mockRepo := new(MockRepository)

	// Arrange
	service := withdraw.NewService(mockRepo, mockValidator)

	// Act
	err := service.Withdraw("", "Ethereum", 100.50)

	// Assert
	assert.Error(t, err)
	assert.Equal(t, "invalid input parameters", err.Error())
}

func TestWithdrawService_ValidationFailed(t *testing.T) {
	mockValidator := new(MockValidationAdapter)
	mockRepo := new(MockRepository)

	// Arrange
	service := withdraw.NewService(mockRepo, mockValidator)
	mockValidator.On("One", "0x123abc456def", "Ethereum").Return(errors.New("wallet validation failed"))

	// Act
	err := service.Withdraw("0x123abc456def", "Ethereum", 100.50)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "wallet validation failed")
}

func TestWithdrawService_RepositoryError(t *testing.T) {
	mockValidator := new(MockValidationAdapter)
	mockRepo := new(MockRepository)

	// Arrange
	service := withdraw.NewService(mockRepo, mockValidator)
	mockValidator.On("One", "0x123abc456def", "Ethereum").Return(nil)
	mockRepo.On("Withdraw", "0x123abc456def", "Ethereum", 100.50).Return(errors.New("insufficient balance"))

	// Act
	err := service.Withdraw("0x123abc456def", "Ethereum", 100.50)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "withdraw transaction failed")
}
