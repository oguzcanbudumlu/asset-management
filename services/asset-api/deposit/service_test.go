package deposit_test

import (
	"asset-management/services/asset-api/deposit"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repository and validation adapter
type mockRepository struct{ mock.Mock }
type mockValidationAdapter struct{ mock.Mock }

func (m *mockValidationAdapter) Both(from, to, network string) error {
	args := m.Called(from, to, network)
	return args.Error(0)
}

func (m *mockRepository) Deposit(walletAddress, network string, amount float64) (float64, error) {
	args := m.Called(walletAddress, network, amount)
	return args.Get(0).(float64), args.Error(1)
}

func (m *mockValidationAdapter) One(walletAddress, network string) error {
	args := m.Called(walletAddress, network)
	return args.Error(0)
}

func TestDepositService_ValidDeposit(t *testing.T) {
	adapter := new(mockValidationAdapter)
	repo := new(mockRepository)

	adapter.On("One", "0x123abc456def", "Ethereum").Return(nil)
	repo.On("Deposit", "0x123abc456def", "Ethereum", 100.50).Return(1500.75, nil)

	service := deposit.NewService(adapter, repo)
	newBalance, err := service.Deposit("0x123abc456def", "Ethereum", 100.50)

	assert.NoError(t, err)
	assert.Equal(t, 1500.75, newBalance)
}

func TestDepositService_InvalidInput(t *testing.T) {
	adapter := new(mockValidationAdapter)
	repo := new(mockRepository)

	service := deposit.NewService(adapter, repo)
	newBalance, err := service.Deposit("", "Ethereum", 100.50)

	assert.Error(t, err)
	assert.Equal(t, 0.0, newBalance)
}

func TestDepositService_ValidationError(t *testing.T) {
	adapter := new(mockValidationAdapter)
	repo := new(mockRepository)

	adapter.On("One", "0x123abc456def", "Ethereum").Return(errors.New("wallet validation failed"))

	service := deposit.NewService(adapter, repo)
	newBalance, err := service.Deposit("0x123abc456def", "Ethereum", 100.50)

	assert.Error(t, err)
	assert.Equal(t, 0.0, newBalance)
}

func TestDepositService_RepositoryError(t *testing.T) {
	adapter := new(mockValidationAdapter)
	repo := new(mockRepository)

	adapter.On("One", "0x123abc456def", "Ethereum").Return(nil)
	repo.On("Deposit", "0x123abc456def", "Ethereum", 100.50).Return(0.0, errors.New("repository error"))

	service := deposit.NewService(adapter, repo)
	newBalance, err := service.Deposit("0x123abc456def", "Ethereum", 100.50)

	assert.Error(t, err)
	assert.Equal(t, 0.0, newBalance)
}
