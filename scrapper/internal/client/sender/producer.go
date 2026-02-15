package sender

import (
	"context"
	"encoding/json"
	"fmt"
	"scrapper/internal/config"
	"scrapper/internal/domain"

	"github.com/segmentio/kafka-go"
)

type Producer struct {
	kafkaWriter *kafka.Writer
}

func NewProducer(cfg *config.Config) *Producer {
	w := &kafka.Writer{
		Addr:         kafka.TCP(cfg.Kafka.Addr),
		Topic:        cfg.Kafka.Topic,
		Balancer:     &kafka.LeastBytes{},
		BatchTimeout: cfg.Kafka.Timeout,
		MaxAttempts:  cfg.Kafka.Retry,
	}
	return &Producer{
		kafkaWriter: w,
	}
}

func (p *Producer) SendUpdate(ctx context.Context, update *domain.Response) error {
	value, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("json marshal: %s", err)
	}

	err = p.kafkaWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(update.Link.URL),
		Value: value,
	})
	if err != nil {
		return fmt.Errorf("failed to send message to sender: %s", err)
	}

	return nil
}

func (p *Producer) Close() error {
	return p.kafkaWriter.Close()
}
