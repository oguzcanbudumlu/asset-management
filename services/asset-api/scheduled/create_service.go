package scheduled

import (
	"asset-management/internal/schedule"
	"asset-management/services/asset-api/wallet"
	"errors"
	"time"
)

type CreateService interface {
	Create(fromWallet, toWallet, network string, amount float64, scheduledTime time.Time) (int, error)
}

type createService struct {
	repo            CreateRepository
	walletValidator wallet.ValidationAdapter
}

func NewCreateService(repo CreateRepository, wv wallet.ValidationAdapter) CreateService {
	return &createService{repo: repo, walletValidator: wv}
}

func (s *createService) Create(fromWallet, toWallet, network string, amount float64, scheduledTime time.Time) (int, error) {
	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}

	if err := s.walletValidator.Both(fromWallet, toWallet, network); err != nil {
		return 0, err
	}

	tx := &schedule.ScheduledTransaction{
		FromWallet:    fromWallet,
		ToWallet:      toWallet,
		Network:       network,
		Amount:        amount,
		ScheduledTime: scheduledTime,
		Status:        schedule.StatusPending,
	}
	return s.repo.Create(tx)
}
