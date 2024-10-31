package main

import (
	"asset-management/internal/schedule"
	"asset-management/internal/schedule/scheduled_process"
	"asset-management/pkg/database"
	"asset-management/pkg/logger"
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"os"
)

func main() {
	logger.InitLogger(zerolog.InfoLevel)
	kafkaBroker := os.Getenv("KAFKA_BROKER")
	kafkaTopic := os.Getenv("KAFKA_TOPIC")
	kafkaGroupId := os.Getenv("KAFKA_GROUP_ID")

	if kafkaBroker == "" || kafkaTopic == "" || kafkaGroupId == "" {
		log.Fatal().Msg("Kafka configuration environment variables (KAFKA_BROKER, KAFKA_TOPIC, KAFKA_GROUP_ID) must be set")
		return
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{kafkaBroker}, // Kafka broker address
		Topic:   kafkaTopic,            // Topic to consume from
		GroupID: kafkaGroupId,          // Consumer group ID
	})

	defer func() {
		if err := reader.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close Kafka reader")
		}
	}()

	if _, err := reader.FetchMessage(context.Background()); err != nil {
		log.Error().Err(err).Msg("Failed to connect to Kafka")
		return
	}
	log.Info().
		Str("broker", kafkaBroker).
		Str("topic", kafkaTopic).
		Str("group_id", kafkaGroupId).
		Msg("Successfully connected to Kafka")

	db, err := database.NewDatabaseRaw(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize database")
		return
	}
	processRepo := scheduled_process.NewProcessRepository(db.Conn)
	processServ := scheduled_process.NewProcessService(processRepo)

	// Start consuming messages
	for {
		// Read messages from Kafka
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Error reading message")
			break
		}

		// Parse message to ScheduleTransaction
		var transaction schedule.ScheduleTransaction
		err = json.Unmarshal(msg.Value, &transaction)
		if err != nil {
			log.Error().Err(err).Msg("Error parsing message")
			continue
		}

		// Process the transaction
		if err := processServ.Process(transaction.ID); err != nil {
			log.Error().
				Err(err).
				Interface("transaction", transaction).
				Msg("Error processing transaction")
			continue
		}

		// Log the successful transaction details
		log.Info().
			Interface("transaction", transaction).
			Msg("Consumed transaction")
	}
}
