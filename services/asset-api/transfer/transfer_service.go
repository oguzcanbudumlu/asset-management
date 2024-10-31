package transfer

import (
	"asset-management/services/asset-api/wallet"
)

type transferService struct {
	transferRepository TransferRepository
	walletValidator    wallet.ValidationAdapter
}

type TransferService interface {
	Transfer(from, to, network string, amount float64) error
}

func NewTransferService(tr TransferRepository, wv wallet.ValidationAdapter) TransferService {
	return &transferService{transferRepository: tr, walletValidator: wv}
}

func (s *transferService) Transfer(from, to, network string, amount float64) error {
	if err := s.walletValidator.Both(from, to, network); err != nil {
		return err
	}

	return s.transferRepository.Transfer(from, to, network, amount)
}
