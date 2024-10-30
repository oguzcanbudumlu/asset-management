package main

import (
	"asset-management/internal/schedule"
	"asset-management/pkg/kafka"
)

type service struct {
	fetchService schedule.NextService
	producer     kafka.Producer
}

type Service interface {
	Execute() (int, error)
}

func (s *service) Execute() (int, error) {
	//transactions, err := s.fetchService.GetNextMinuteTransactions()
	//if err != nil {
	//	return 0, fmt.Errorf("failed to execute service: %v", err)
	//}
	//
	return 0, nil
}
