package usecase

import (
	"context"
	"errors"
)

func (u *UseCase) CreateChat(ctx context.Context, chatID int64) error {
	if chatID == 0 {
		return errors.New("invalid chat id")
	}

	err := u.db.CreateChat(ctx, chatID)
	if err != nil {
		u.log.Error("failed to create chat", "chatID", chatID, "error", err)
		return err
	}

	return nil
}
