package usecase

import (
	"context"
	"log/slog"
	"scrapper/internal/config"
	"scrapper/internal/domain"
)

type Postgres interface {
	DeleteUserLink(ctx context.Context, chatID int64, alias string) error
	DeleteLink(ctx context.Context, link *domain.Link) error
	AddLink(ctx context.Context, link *domain.Link) error
	AddUserLink(ctx context.Context, chatID int64, link *domain.Link) error
	GetLinksByChatID(ctx context.Context, chatID int64) ([]domain.Link, error)
	GetUserLinksByTag(ctx context.Context, chatID int64, tags string) ([]*domain.Link, error)
	GetLinks(ctx context.Context, limit, offset uint64) ([]domain.Link, error)
	GetLinkByURL(ctx context.Context, url string) (*domain.Link, error)
	IsLinkExists(ctx context.Context, url string) (bool, error)
	IsUserLinkExists(ctx context.Context, alias string, chatID int64) (bool, error)
	DeleteChat(ctx context.Context, chatID int64) error
	CreateChat(ctx context.Context, chatID int64) error
}

type UseCase struct {
	db  Postgres
	log *slog.Logger
	ctx context.Context
	cfg *config.Config
}

func NewUseCase(db Postgres, log *slog.Logger, ctx context.Context, cfg *config.Config) *UseCase {
	return &UseCase{
		db:  db,
		log: log,
		ctx: ctx,
		cfg: cfg,
	}
}
