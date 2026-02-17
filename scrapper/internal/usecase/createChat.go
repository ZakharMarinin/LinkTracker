package usecase

import (
	"context"
	"errors"
)

func (u *UseCase) CreateChat(ctx context.Context, chatID int64) error {
	if chatID == 0 {
		return errors.New("invalid chat id")
	}

	err := u.DB.CreateChat(ctx, chatID)
	if err != nil {
		u.Log.Error("failed to create chat", "chatID", chatID, "error", err)
		return err
	}

	return nil
}
