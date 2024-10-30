package withdraw

import (
	"asset-management/services/asset-api/wallet"
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
)

type withdrawService struct {
	withdrawRepository WithdrawRepository
	walletValidator    wallet.ValidationAdapter
}

type WithdrawService interface {
	Withdraw(walletAddress, network string, amount decimal.Decimal) (decimal.Decimal, error)
}

func NewWithdrawService(wr WithdrawRepository, va wallet.ValidationAdapter) WithdrawService {
	return &withdrawService{withdrawRepository: wr, walletValidator: va}
}

func (s *withdrawService) Withdraw(walletAddress, network string, amount decimal.Decimal) (decimal.Decimal, error) {
	if walletAddress == "" || network == "" || amount.LessThanOrEqual(decimal.Zero) {
		return decimal.Zero, errors.New("invalid input parameters")
	}

	// Validate the wallet
	err := s.walletValidator.One(walletAddress, network)
	if err != nil {
		return decimal.Zero, fmt.Errorf("wallet validation failed: %w", err)
	}

	// Perform the withdrawal
	newBalance, repoErr := s.withdrawRepository.Withdraw(walletAddress, network, amount)
	if repoErr != nil {
		return decimal.Zero, fmt.Errorf("withdraw transaction failed: %w", repoErr)
	}

	return newBalance, nil
}
