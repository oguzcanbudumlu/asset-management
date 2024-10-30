package kafka

import (
	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"time"
)

type Producer struct {
	Writer *kafka.Writer
}

func NewProducer(brokerAddress, topic string) *Producer {
	writer := &kafka.Writer{
		Addr:      kafka.TCP(brokerAddress),
		Topic:     topic,
		Balancer:  &kafka.LeastBytes{},
		Transport: &kafka.Transport{ClientID: "json-producer", DialTimeout: 10 * time.Second},
	}
	return &Producer{Writer: writer}
}

// Close shuts down the Kafka Writer connection
func (p *Producer) Close() error {
	if err := p.Writer.Close(); err != nil {
		log.Error().Err(err).Msg("Failed to close Kafka producer")
		return err
	}
	log.Info().Msg("Kafka producer closed successfully")
	return nil
}
