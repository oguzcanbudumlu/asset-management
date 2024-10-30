package publisher

import (
	"asset-management/internal/schedule"
	"asset-management/pkg/kafka"
)

type service struct {
	nextService schedule.NextService
	producer    *kafka.Producer
}

type Service interface {
	TriggerPublisher() (int, error)
}

func NewService(nextService schedule.NextService, producer *kafka.Producer) Service {
	return &service{
		nextService: nextService,
		producer:    producer,
	}
}

func (s *service) TriggerPublisher() (int, error) {
	//transactions, err := s.nextService.GetNextMinuteTransactions()
	return 0, nil
	//if err != nil {
	//	return 0, fmt.Errorf("failed to retrieve transactions: %v", err)
	//}
	//
	//// Map transactions to events for Kafka production
	//events := MapScheduleTransactionsToEvents(transactions)
	//
	//for _, event := range events {
	//	key, value := MapScheduleTransactionsToEvents(event) // Assume MapScheduleTransactionToKeyValue exists
	//	keyBytes, _ := json.Marshal(key)
	//	valueBytes, _ := json.Marshal(value)
	//
	//	messages = append(messages, kafka.Message{
	//		Key:   keyBytes,
	//		Value: valueBytes,
	//	})
	//}
	//
	//// Produce the events to Kafka
	//if err := s.producer.ProduceMessages(messages); err != nil {
	//	return 0, fmt.Errorf("failed to produce messages to Kafka: %v", err)
	//}
	//
	//return len(messages), nil
}
