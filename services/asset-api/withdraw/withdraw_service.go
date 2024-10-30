package withdraw

import (
	"asset-management/internal/wallet"
	"errors"
	"fmt"
)

type withdrawService struct {
	withdrawRepository WithdrawRepository
	walletValidator    wallet.ValidationAdapter
}

type WithdrawService interface {
	Withdraw(walletAddress, network string, amount float64) error
}

func NewWithdrawService(wr WithdrawRepository, va wallet.ValidationAdapter) WithdrawService {
	return &withdrawService{withdrawRepository: wr, walletValidator: va}
}

func (s *withdrawService) Withdraw(walletAddress, network string, amount float64) error {
	if walletAddress == "" || network == "" || amount <= 0 {
		return errors.New("invalid input parameters")
	}

	err := s.walletValidator.One(walletAddress, network)

	if err != nil {
		return fmt.Errorf("wallet validation failed: %w", err)
	}

	repoErr := s.withdrawRepository.Withdraw(walletAddress, network, amount)
	if repoErr != nil {
		return fmt.Errorf("withdraw transaction failed: %w", err)
	}

	return nil
}
