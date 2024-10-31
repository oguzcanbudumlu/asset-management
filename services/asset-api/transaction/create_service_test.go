package transaction

import (
	"asset-management/internal/schedule"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type MockCreateRepository struct {
	mock.Mock
}

func (m *MockCreateRepository) Create(tx *schedule.ScheduleTransaction) (int, error) {
	args := m.Called(tx)
	return args.Int(0), args.Error(1)
}

type MockValidationAdapter struct {
	mock.Mock
}

func (m *MockValidationAdapter) Both(from, to, network string) error {
	args := m.Called(from, to, network)
	return args.Error(0)
}

func (m *MockValidationAdapter) One(wallet, network string) error {
	args := m.Called(wallet, network)
	return args.Error(0)
}

func TestCreateService_Success(t *testing.T) {
	mockRepo := new(MockCreateRepository)
	mockValidator := new(MockValidationAdapter)
	service := NewCreateService(mockRepo, mockValidator)

	mockValidator.On("Both", "wallet123", "wallet456", "mainnet").Return(nil)
	mockRepo.On("Create", mock.Anything).Return(123, nil)

	id, err := service.Create("wallet123", "wallet456", "mainnet", 100.50, time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 123, id)

	mockValidator.AssertExpectations(t)
	mockRepo.AssertExpectations(t)
}

func TestCreateService_ValidationError(t *testing.T) {
	mockRepo := new(MockCreateRepository)
	mockValidator := new(MockValidationAdapter)
	service := NewCreateService(mockRepo, mockValidator)

	mockValidator.On("Both", "wallet123", "wallet456", "mainnet").Return(errors.New("validation failed"))

	id, err := service.Create("wallet123", "wallet456", "mainnet", 100.50, time.Now())
	assert.Error(t, err)
	assert.Equal(t, "validation failed", err.Error())
	assert.Equal(t, 0, id)

	mockValidator.AssertExpectations(t)
}

func TestCreateService_InvalidAmount(t *testing.T) {
	mockRepo := new(MockCreateRepository)
	mockValidator := new(MockValidationAdapter)
	service := NewCreateService(mockRepo, mockValidator)

	id, err := service.Create("wallet123", "wallet456", "mainnet", 0, time.Now())
	assert.Error(t, err)
	assert.Equal(t, "amount must be greater than zero", err.Error())
	assert.Equal(t, 0, id)
}
