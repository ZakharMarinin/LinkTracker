package usecase

import (
	"context"
)

func (u *UseCase) CreateChat(ctx context.Context, chatID int64) error {
	err := u.db.CreateChat(ctx, chatID)
	if err != nil {
		u.log.Error("failed to create chat", "chatID", chatID, "error", err)
		return err
	}

	return nil
}
