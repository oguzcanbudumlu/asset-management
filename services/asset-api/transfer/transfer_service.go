package transfer

import (
	"asset-management/internal/common"
	"errors"
	"golang.org/x/sync/errgroup"
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
	if err := s.validateWallets(from, to, network); err != nil {
		return err
	}

	return s.transferRepository.Transfer(from, to, network, amount)
}

func (s *transferService) validateWallets(from, to, network string) error {
	var g errgroup.Group

	g.Go(func() error {
		if err := s.walletValidator.ValidateWallet(from, network); err != nil {
			return errors.New("source wallet validation failed: " + err.Error())
		}
		return nil
	})

	g.Go(func() error {
		if err := s.walletValidator.ValidateWallet(to, network); err != nil {
			return errors.New("destination wallet validation failed: " + err.Error())
		}
		return nil
	})

	return g.Wait()
}
