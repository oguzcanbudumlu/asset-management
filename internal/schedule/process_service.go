package schedule

import (
	"fmt"
)

type ProcessService interface {
	Process(scheduledTransactionID int) error
}

type processService struct {
	repo ProcessRepository
}

func NewProcessService(repo ProcessRepository) ProcessService {
	return &processService{repo: repo}
}

func (s *processService) Process(scheduledTransactionID int) error {
	err := s.repo.Process(scheduledTransactionID)
	if err != nil {
		return fmt.Errorf("failed to process transaction: %w", err)
	}

	return nil
}
