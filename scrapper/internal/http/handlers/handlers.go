package handlers

import (
	"context"
	"log/slog"
	"scrapper/internal/domain"
)

type UseCase interface {
	CreateChat(ctx context.Context, chatID int64) error
	DeleteChat(ctx context.Context, chatID int64) error
	AddLink(ctx context.Context, chatID int64, url string, desc string) error
	DeleteLink(ctx context.Context, chatID int64, alias string) error
	GetLinks(ctx context.Context, chatID int64) ([]domain.Link, error)
}

type HTTP struct {
	useCase UseCase
	log     *slog.Logger
}

func NewHTTP(useCase UseCase, log *slog.Logger) *HTTP {
	return &HTTP{useCase: useCase, log: log}
}
