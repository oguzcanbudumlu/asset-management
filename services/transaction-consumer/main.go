package main

import (
	"asset-management/pkg/logger"
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"os"
)

func main() {
	logger.InitLogger(zerolog.InfoLevel)
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	log.Info().Msg(kafkaBroker)
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker}, // Kafka broker address
		Topic:   "test-topic",          // Topic to consume from
		GroupID: "consumer-group-1",    // Consumer group ID
	})

	defer reader.Close()

	log.Info().Msg("Consuming messages from Kafka...")

	for {
		// Read message
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Error reading message")
			continue
		}

		// Log message details
		log.Info().
			Str("topic", msg.Topic).
			Int("partition", msg.Partition).
			Int64("offset", msg.Offset).
			Str("key", string(msg.Key)).
			Str("value", string(msg.Value)).
			Msg("Received new message")
	}
}
