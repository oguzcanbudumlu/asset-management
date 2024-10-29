package withdraw

import (
	"asset-management/internal/common"
	"errors"
	"fmt"
)

type withdrawService struct {
	withdrawRepository WithdrawRepository
	walletValidator    common.WalletValidationAdapter
}

type WithdrawService interface {
	Withdraw(walletAddress, network string, amount float64) error
}

func NewWithdrawService(wr WithdrawRepository, va common.WalletValidationAdapter) WithdrawService {
	return &withdrawService{withdrawRepository: wr, walletValidator: va}
}

func (s *withdrawService) Withdraw(walletAddress, network string, amount float64) error {
	if walletAddress == "" || network == "" || amount <= 0 {
		return errors.New("invalid input parameters")
	}

	err := s.walletValidator.ValidateWallet(walletAddress, network)

	if err != nil {
		return fmt.Errorf("wallet validation failed: %w", err)
	}

	repoErr := s.withdrawRepository.Withdraw(walletAddress, network, amount)
	if repoErr != nil {
		return fmt.Errorf("withdraw transaction failed: %w", err)
	}

	return nil
}
