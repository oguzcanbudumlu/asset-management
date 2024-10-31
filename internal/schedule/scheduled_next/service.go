package scheduled_next

import "asset-management/internal/schedule"

type NextService interface {
	GetNextMinuteTransactions() ([]schedule.ScheduleTransaction, error)
}

type nextService struct {
	repo NextRepository
}

func NewNextService(repo NextRepository) NextService {
	return &nextService{repo: repo}
}

func (s *nextService) GetNextMinuteTransactions() ([]schedule.ScheduleTransaction, error) {
	return s.repo.GetNextMinuteTransactions()
}
