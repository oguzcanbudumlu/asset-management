package e2e

import "github.com/stretchr/testify/mock"

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
