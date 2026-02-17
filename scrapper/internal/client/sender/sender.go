package sender

import (
	"context"
	"fmt"
	"log/slog"
	"scrapper/internal/config"
	"scrapper/internal/domain"
)

type Sender interface {
	SendUpdate(ctx context.Context, link *domain.Response) error
}

type TypeOfSender struct {
	Primary   Sender
	Secondary Sender
	Logger    *slog.Logger
}

func New(logger *slog.Logger, cfg *config.Config) (*TypeOfSender, error) {
	httpSender := NewTGClient(logger, cfg)
	kafkaSender := NewProducer(cfg)

	switch cfg.TransportType {
	case "kafka":
		return &TypeOfSender{
			Primary:   kafkaSender,
			Secondary: httpSender,
			Logger:    logger,
		}, nil
	case "http":
		return &TypeOfSender{
			Primary:   httpSender,
			Secondary: kafkaSender,
			Logger:    logger,
		}, nil
	default:
		return nil, fmt.Errorf("unknown message type: %s", cfg.TransportType)
	}
}

func (s *TypeOfSender) SendUpdate(ctx context.Context, link *domain.Response) error {
	err := s.Primary.SendUpdate(ctx, link)
	if err == nil {
		return nil
	}

	s.Logger.Warn("primary transport method failed", "error", err)

	err = s.Secondary.SendUpdate(ctx, link)
	if err != nil {
		s.Logger.Error("secondary transport method failed", "error", err)
	}

	return nil
}
