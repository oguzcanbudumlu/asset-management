package main

import (
	"asset-management/pkg/logger"
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"os"
	"time"
)

type Event struct {
	UserID    string `json:"user_id"`
	Action    string `json:"action"`
	Timestamp int64  `json:"timestamp"`
}

func main() {

	logger.InitLogger(zerolog.InfoLevel)
	brokerAddress := os.Getenv("KAFKA_BROKER") // Replace with your Kafka broker address
	topic := "test-topic"                      // Replace with your topic name

	// Set up Kafka transport
	transport := &kafka.Transport{
		ClientID:    "json-producer",
		DialTimeout: 10 * time.Second,
	}

	// Create a Kafka writer
	writer := &kafka.Writer{
		Addr:      kafka.TCP(brokerAddress),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		Transport: transport,
	}

	// Create an event
	event := Event{
		UserID:    "12345",
		Action:    "login",
		Timestamp: time.Now().Unix(),
	}

	// Serialize the event to JSON
	messageBytes, err := json.Marshal(event)
	if err != nil {
		log.Error().Err(err).Msg("Failed to serialize event to JSON")
		return
	}

	// Create a Kafka message with the JSON payload
	// Using the UserID as a string key
	message := kafka.Message{
		Key:   []byte(event.UserID), // Convert string key to []byte
		Value: messageBytes,
	}

	// Log before sending the message
	log.Info().
		Str("topic", topic).
		Str("broker", brokerAddress).
		Str("event", string(messageBytes)).
		Msg("Producing JSON message to Kafka")

	// Produce (send) the message to Kafka
	if err := writer.WriteMessages(context.Background(), message); err != nil {
		log.Error().
			Err(err).
			Str("topic", topic).
			Msg("Failed to produce JSON message to Kafka")
	} else {
		log.Info().
			Str("topic", topic).
			Msg("JSON message produced successfully to Kafka")
	}

	// Close the writer explicitly
	if err := writer.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close Kafka writer")
	} else {
		log.Info().Msg("Kafka writer closed successfully")
	}
}
