package withdraw

import (
	"asset-management/services/asset-api/wallet"
	"errors"
	"fmt"
)

type service struct {
	withdrawRepository Repository
	walletValidator    wallet.ValidationAdapter
}

type Service interface {
	Withdraw(walletAddress, network string, amount float64) error
}

func NewService(wr Repository, va wallet.ValidationAdapter) Service {
	return &service{withdrawRepository: wr, walletValidator: va}
}

func (s *service) Withdraw(walletAddress, network string, amount float64) error {
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
