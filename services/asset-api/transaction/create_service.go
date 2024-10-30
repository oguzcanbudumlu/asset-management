package transaction

import (
	"asset-management/internal/common"
	"asset-management/internal/schedule"
	"errors"
	"time"
)

type CreateService interface {
	Create(fromWallet, toWallet, network string, amount float64, scheduledTime time.Time) (int, error)
}

type createService struct {
	repo            CreateRepository
	walletValidator common.WalletValidationAdapter
}

func NewCreateService(repo CreateRepository, wv common.WalletValidationAdapter) CreateService {
	return &createService{repo: repo, walletValidator: wv}
}

func (s *createService) Create(fromWallet, toWallet, network string, amount float64, scheduledTime time.Time) (int, error) {
	if err := s.walletValidator.ValidateBoth(fromWallet, toWallet, network); err != nil {
		return 0, err
	}

	if amount <= 0 {
		return 0, errors.New("amount must be greater than zero")
	}
	tx := &schedule.ScheduleTransaction{
		FromWallet:    fromWallet,
		ToWallet:      toWallet,
		Network:       network,
		Amount:        amount,
		ScheduledTime: scheduledTime,
		Status:        schedule.StatusPending,
	}
	return s.repo.Create(tx)
}
