package main

import (
	"errors"
	"fmt"
)

type depositService struct {
	depositRepository DepositRepository
	validationAdapter WalletValidationAdapter
}

type DepositService interface {
	Deposit(walletAddress, network string, amount float64) (string, float64, error)
}

func NewDepositService(adapter WalletValidationAdapter, depositRepository DepositRepository) DepositService {
	return &depositService{validationAdapter: adapter, depositRepository: depositRepository}
}
func (s *depositService) Deposit(walletAddress, network string, amount float64) (string, float64, error) {
	// Validate input
	if walletAddress == "" || network == "" || amount <= 0 {
		return "", 0, errors.New("invalid input parameters")
	}

	// Validate the wallet using WalletValidationAdapter
	isValid, err := s.validationAdapter.ValidateWallet(walletAddress, network)
	if err != nil {
		return "", 0, fmt.Errorf("wallet validation failed: %w", err)
	}
	if !isValid {
		return "", 0, errors.New("wallet is invalid or inactive")
	}

	// Perform the deposit transaction
	newBalance, err := s.depositRepository.Deposit(walletAddress, network, amount)
	if err != nil {
		return "", 0, fmt.Errorf("deposit transaction failed: %w", err)
	}

	// Generate a mock transaction ID for the response
	transactionID := "txn_" + walletAddress[:6] + "123"

	return transactionID, newBalance, nil
}
