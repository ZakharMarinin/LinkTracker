package usecase

import (
	"context"
	"linktracker/internal/domain"
	"linktracker/internal/storage"
	"log/slog"
)

type Storage interface {
	SetTempUserState(ctx context.Context, userInfo *domain.UserStateInfo) error
	GetTempUserState(ctx context.Context, userID int64) (*domain.UserStateInfo, error)
	SaveTempUserLinks(ctx context.Context, tempUserLinks *storage.TempUserLinks) error
	GetTempUserLinks(ctx context.Context, userID int64) (*storage.TempUserLinks, error)
}

type ScrapperClient interface {
	CreateChat(ctx context.Context, chatID int64) error
	AddLink(ctx context.Context, chatID int64, link domain.Link) error
	DeleteChat(ctx context.Context, chatID int64) error
	DeleteLink(ctx context.Context, chatID int64, alias string) error
	GetLinks(ctx context.Context, chatID int64) ([]*domain.Link, error)
	GetFilteredLinks(ctx context.Context, chatID int64, tag string) ([]*domain.Link, error)
}

type UseCase struct {
	log            *slog.Logger
	ScrapperClient ScrapperClient
	Storage        Storage
}

func New(log *slog.Logger, ScrapperClient ScrapperClient, Storage Storage) *UseCase {
	return &UseCase{
		log:            log,
		ScrapperClient: ScrapperClient,
		Storage:        Storage,
	}
}

func (u *UseCase) ChangeUserState(ctx context.Context, userInfo *domain.UserStateInfo, state string) error {
	userInfo.State = state
	err := u.Storage.SetTempUserState(ctx, userInfo)
	if err != nil {
		u.log.Error("ChangeState: Error setting user state", "error", err)
		return err
	}
	return nil
}

func (u *UseCase) GetUserState(ctx context.Context, userID int64) (*domain.UserStateInfo, error) {
	userInfo, err := u.Storage.GetTempUserState(ctx, userID)
	if err != nil {
		u.log.Error("GetUserState: Error getting user state", "error", err)
		return nil, err
	}

	return userInfo, nil
}
