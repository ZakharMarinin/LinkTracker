package kafka

import (
	"context"
	"fmt"
	"linktracker/internal/config"

	"github.com/segmentio/kafka-go"
)

type Consumer struct {
	kafkaReader *kafka.Reader
}

func NewConsumer(cfg *config.Config) *Consumer {
	return &Consumer{
		kafkaReader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{cfg.BotClients.Kafka.Addr},
			Topic:          cfg.BotClients.Kafka.Topic,
			GroupID:        cfg.BotClients.Kafka.GroupID,
			StartOffset:    kafka.FirstOffset,
			MaxWait:        cfg.BotClients.Kafka.Timeout,
			MaxAttempts:    cfg.BotClients.Kafka.Retry,
			CommitInterval: 0,
		}),
	}
}

func (c *Consumer) ReadMessage(ctx context.Context) (*kafka.Message, error) {
	msg, err := c.kafkaReader.FetchMessage(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to read Kafka message: %s", err)
	}

	return &msg, nil
}

func (c *Consumer) CommitMessage(ctx context.Context, msg *kafka.Message) error {
	err := c.kafkaReader.CommitMessages(ctx, *msg)
	if err != nil {
		return fmt.Errorf("failed to commit Kafka message: %s", err)
	}

	return nil
}

func (c *Consumer) Close() error {
	return c.kafkaReader.Close()
}
