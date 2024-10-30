package deposit

import (
	"asset-management/internal/wallet"
	"errors"
	"fmt"
)

type depositService struct {
	depositRepository DepositRepository
	validationAdapter wallet.ValidationAdapter
}

type DepositService interface {
	Deposit(walletAddress, network string, amount float64) (float64, error)
}

func NewDepositService(adapter wallet.ValidationAdapter, depositRepository DepositRepository) DepositService {
	return &depositService{validationAdapter: adapter, depositRepository: depositRepository}
}
func (s *depositService) Deposit(walletAddress, network string, amount float64) (float64, error) {
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
