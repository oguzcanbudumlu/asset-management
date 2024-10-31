package scheduled_process

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type MockProcessRepository struct {
	mock.Mock
}

func (m *MockProcessRepository) Process(scheduledTransactionID int) error {
	args := m.Called(scheduledTransactionID)
	return args.Error(0)
}

func TestProcessService_Process_Success(t *testing.T) {
	mockRepo := new(MockProcessRepository)
	service := NewProcessService(mockRepo)

	// Mock successful repository response
	mockRepo.On("Process", 123).Return(nil)

	err := service.Process(123)

	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestProcessService_Process_Failure(t *testing.T) {
	mockRepo := new(MockProcessRepository)
	service := NewProcessService(mockRepo)

	// Mock repository error
	mockRepo.On("Process", 123).Return(errors.New("repository error"))

	err := service.Process(123)

	// Assertions
	assert.Error(t, err)
	assert.Equal(t, "failed to process transaction: repository error", err.Error())
	mockRepo.AssertExpectations(t)
}
