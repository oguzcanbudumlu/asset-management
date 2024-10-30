package schedule

type NextService interface {
	GetNextMinuteTransactions() ([]ScheduleTransaction, error)
}

type nextService struct {
	repo NextRepository
}

func NewNextService(repo NextRepository) NextService {
	return &nextService{repo: repo}
}

func (s *nextService) GetNextMinuteTransactions() ([]ScheduleTransaction, error) {
	return s.repo.GetNextMinuteTransactions()
}
