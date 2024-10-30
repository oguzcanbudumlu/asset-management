package transaction

import (
	"asset-management/internal/schedule"
	"asset-management/services/asset-api/wallet"
	"errors"
	"github.com/shopspring/decimal"
	"time"
)

type CreateService interface {
	Create(fromWallet, toWallet, network string, amount decimal.Decimal, scheduledTime time.Time) (int, error)
}

type createService struct {
	repo            CreateRepository
	walletValidator wallet.ValidationAdapter
}

func NewCreateService(repo CreateRepository, wv wallet.ValidationAdapter) CreateService {
	return &createService{repo: repo, walletValidator: wv}
}

func (s *createService) Create(fromWallet, toWallet, network string, amount decimal.Decimal, scheduledTime time.Time) (int, error) {
	if err := s.walletValidator.Both(fromWallet, toWallet, network); err != nil {
		return 0, err
	}

	if amount.LessThanOrEqual(decimal.Zero) {
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
