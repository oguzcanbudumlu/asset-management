package deposit

import (
	"asset-management/services/asset-api/wallet"
	"errors"
	"fmt"
)

type service struct {
	depositRepository Repository
	validationAdapter wallet.ValidationAdapter
}

type Service interface {
	Deposit(walletAddress, network string, amount float64) (float64, error)
}

func NewService(adapter wallet.ValidationAdapter, depositRepository Repository) Service {
	return &service{validationAdapter: adapter, depositRepository: depositRepository}
}
func (s *service) Deposit(walletAddress, network string, amount float64) (float64, error) {
	// Validate input
	if walletAddress == "" || network == "" || amount <= 0 {
		return 0, errors.New("invalid input parameters")
	}

	err := s.validationAdapter.One(walletAddress, network)
	if err != nil {
		return 0, fmt.Errorf("wallet validation failed: %w", err)
	}

	// Perform the deposit transaction
	newBalance, err := s.depositRepository.Deposit(walletAddress, network, amount)
	if err != nil {
		return 0, fmt.Errorf("deposit transaction failed: %w", err)
	}

	return newBalance, nil
}
