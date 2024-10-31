package scheduled_next_test

import (
	"asset-management/internal/schedule"
	"asset-management/internal/schedule/scheduled_next"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockNextRepository struct {
	mock.Mock
}

func (m *MockNextRepository) GetNextMinuteTransactions() ([]schedule.ScheduledTransaction, error) {
	args := m.Called()
	if transactions, ok := args.Get(0).([]schedule.ScheduledTransaction); ok {
		return transactions, args.Error(1)
	}
	return nil, args.Error(1)
}

func TestNextService_GetNextMinuteTransactions_Success(t *testing.T) {
	mockRepo := new(MockNextRepository)
	service := scheduled_next.NewNextService(mockRepo)

	transactions := []schedule.ScheduledTransaction{
		{ID: 1, FromWallet: "wallet123", ToWallet: "wallet456", Network: "mainnet", Amount: 100.50},
	}
	mockRepo.On("GetNextMinuteTransactions").Return(transactions, nil)

	result, err := service.GetNextMinuteTransactions()

	assert.NoError(t, err)
	assert.Equal(t, transactions, result)
	mockRepo.AssertExpectations(t)
}

func TestNextService_GetNextMinuteTransactions_Error(t *testing.T) {
	mockRepo := new(MockNextRepository)
	service := scheduled_next.NewNextService(mockRepo)

	mockRepo.On("GetNextMinuteTransactions").Return(nil, errors.New("repository error"))

	result, err := service.GetNextMinuteTransactions()

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "repository error", err.Error())
	mockRepo.AssertExpectations(t)
}
