package main

import "errors"

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

	// Validate wallet with the Wallet Management Service
	isValid, err := s.validationAdapter.ValidateWallet(walletAddress, network)
	if err != nil {
		return "", 0, err
	}
	if !isValid {
		return "", 0, errors.New("wallet is invalid or inactive")
	}

	// Retrieve current balance

	return "", 0, nil
}
