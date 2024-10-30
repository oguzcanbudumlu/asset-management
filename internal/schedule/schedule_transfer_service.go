package schedule

//
//import (
//	"asset-management/internal/common"
//	"errors"
//	"fmt"
//	"time"
//)
//
//type ScheduleTransactionService interface {
//	Create(fromWallet, toWallet, network string, amount float64, scheduledTime time.Time) (int, error)
//	GetNextMinuteTransactions() ([]ScheduleTransaction, error)
//	Process(scheduledTransactionID int64) error
//}
//
//type scheduleTransactionService struct {
//	repo            ScheduleTransactionRepository
//	walletValidator common.WalletValidationAdapter
//}
//
//func NewScheduleTransactionService(repo ScheduleTransactionRepository, wv common.WalletValidationAdapter) ScheduleTransactionService {
//	return &scheduleTransactionService{repo: repo, walletValidator: wv}
//}
//
//// ScheduleTransaction creates a new schedule transaction.
//func (s *scheduleTransactionService) Create(fromWallet, toWallet, network string, amount float64, scheduledTime time.Time) (int, error) {
//	if err := s.walletValidator.ValidateBoth(fromWallet, toWallet, network); err != nil {
//		return 0, err
//	}
//
//	if amount <= 0 {
//		return 0, errors.New("amount must be greater than zero")
//	}
//	tx := &ScheduleTransaction{
//		FromWallet:    fromWallet,
//		ToWallet:      toWallet,
//		Network:       network,
//		Amount:        amount,
//		ScheduledTime: scheduledTime,
//		Status:        StatusPending,
//	}
//	return s.repo.Create(tx)
//}
//
//func (s *scheduleTransactionService) GetNextMinuteTransactions() ([]ScheduleTransaction, error) {
//	return s.repo.GetNextMinuteTransactions()
//}
//
//func (s *scheduleTransactionService) Process(scheduledTransactionID int64) error {
//	err := s.repo.Process(scheduledTransactionID)
//	if err != nil {
//		return fmt.Errorf("failed to process transaction: %w", err)
//	}
//
//	return nil
//}
