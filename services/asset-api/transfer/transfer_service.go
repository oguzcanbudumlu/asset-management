package transfer

import (
	"asset-management/internal/common"
)

type transferService struct {
	transferRepository TransferRepository
	walletValidator    common.WalletValidationAdapter
}

type TransferService interface {
	Transfer(from, to, network string, amount float64) error
}

func NewTransferService(tr TransferRepository, wv common.WalletValidationAdapter) TransferService {
	return &transferService{transferRepository: tr, walletValidator: wv}
}

func (s *transferService) Transfer(from, to, network string, amount float64) error {
	if err := s.walletValidator.ValidateBoth(from, to, network); err != nil {
		return err
	}

	return s.transferRepository.Transfer(from, to, network, amount)
}
