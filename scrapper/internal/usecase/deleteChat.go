package usecase

import (
	"context"
	"errors"
)

func (u *UseCase) DeleteChat(ctx context.Context, chatID int64) error {
	if chatID == 0 {
		return errors.New("invalid chat id")
	}

	err := u.db.DeleteChat(ctx, chatID)
	if err != nil {
		u.log.Error("Failed to delete chat", "chatID", chatID)
		return err
	}

	return nil
}
