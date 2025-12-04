package usecase

import "context"

func (u *UseCase) CreateChat(ctx context.Context, chatID int64) error {
	err := u.ScrapperClient.CreateChat(ctx, chatID)
	if err != nil {
		u.log.Error("Failed to create chat", "chatID", chatID, "err", err)
		return err
	}

	return nil
}
