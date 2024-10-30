package deposit

import (
	"asset-management/services/asset-api/wallet"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

type depositService struct {
	depositRepository DepositRepository
	validationAdapter wallet.ValidationAdapter
}

type DepositService interface {
	Deposit(walletAddress, network string, amount decimal.Decimal) (decimal.Decimal, error)
}

func NewDepositService(adapter wallet.ValidationAdapter, depositRepository DepositRepository) DepositService {
	return &depositService{validationAdapter: adapter, depositRepository: depositRepository}
}

func (s *depositService) Deposit(walletAddress, network string, amount decimal.Decimal) (decimal.Decimal, error) {
	// Validate input
	if walletAddress == "" || network == "" || amount.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, errors.New("invalid input parameters")
	}

	// Validate the wallet
	err := s.validationAdapter.One(walletAddress, network)
	if err != nil {
		return decimal.Zero, fmt.Errorf("wallet validation failed: %w", err)
	}

	// Perform the deposit transaction
	newBalance, err := s.depositRepository.Deposit(walletAddress, network, amount)
	if err != nil {
		return decimal.Zero, fmt.Errorf("deposit transaction failed: %w", err)
	}

	return newBalance, nil
}
