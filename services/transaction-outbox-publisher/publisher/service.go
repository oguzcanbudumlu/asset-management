package publisher

import (
	"asset-management/internal/schedule"
	"asset-management/pkg/kafka"
	"context"
	"encoding/json"
	"github.com/rs/zerolog/log"
	kafka2 "github.com/segmentio/kafka-go"
	"os"
	"time"
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
	// Retrieve the topic from environment variables once.
	topic := os.Getenv("KAFKA_TOPIC")

	transactions := mockTransactions()
	//transactions, err := s.nextService.GetNextMinuteTransactions()
	//
	//if err != nil {
	//	return 0, err
	//}

	var messages []kafka2.Message

	// Serialize each transaction to JSON and prepare the Kafka messages.
	for _, transaction := range transactions {
		message, err := json.Marshal(transaction)
		if err != nil {
			log.Error().
				Int("transaction_id", transaction.ID).
				Err(err).
				Msg("Error serializing transaction")
			continue
		}

		// Append the serialized transaction as a Kafka message to the slice.
		messages = append(messages, kafka2.Message{
			Key:   []byte(transaction.FromWallet), // Using FromWallet as the message key
			Value: message,
			Time:  time.Now(),
		})
	}

	// Send all messages in bulk to Kafka.
	writeErr := s.producer.Writer.WriteMessages(context.Background(), messages...)
	if writeErr != nil {
		log.Error().
			Err(writeErr).
			Int("message_count", len(messages)).
			Str("topic", topic).
			Msg("Error writing messages to Kafka")
		return 0, writeErr
	}

	log.Info().
		Int("message_count", len(messages)).
		Str("topic", topic).
		Msg("Successfully sent transactions to Kafka")

	return len(messages), nil
}

func mockTransactions() []schedule.ScheduleTransaction {
	return []schedule.ScheduleTransaction{
		{
			ID:            1,
			FromWallet:    "wallet_ABC123",
			ToWallet:      "wallet_XYZ789",
			Network:       "Ethereum",
			Amount:        250.75,
			ScheduledTime: time.Date(2024, 10, 31, 15, 0, 0, 0, time.UTC),
			Status:        "PENDING",
			CreatedAt:     time.Date(2024, 10, 29, 10, 15, 0, 0, time.UTC),
		},
		{
			ID:            2,
			FromWallet:    "wallet_DEF456",
			ToWallet:      "wallet_UVW123",
			Network:       "Bitcoin",
			Amount:        500.00,
			ScheduledTime: time.Date(2024, 11, 1, 12, 0, 0, 0, time.UTC),
			Status:        "PENDING",
			CreatedAt:     time.Date(2024, 10, 30, 11, 30, 0, 0, time.UTC),
		},
	}
}
