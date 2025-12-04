package usecase

import "context"

func (u *UseCase) DeleteChat(ctx context.Context, chatID int64) error {
	err := u.db.DeleteChat(ctx, chatID)
	if err != nil {
		u.log.Error("Failed to delete chat", "chatID", chatID)
		return err
	}

	return nil
}
