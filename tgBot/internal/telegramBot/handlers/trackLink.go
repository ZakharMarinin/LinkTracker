package handlers

import (
	"context"
	"linktracker/internal/domain"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) TrackLink(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userInfo, err := b.useCase.GetUserState(ctx, c.Sender().ID)
		if err != nil {
			b.log.Warn("TrackLink: Error getting user state. Starting backoff...", "error", err)

			err = b.BackoffStart(ctx, c)
			if err != nil {
				b.log.Error("backoff start error", "error", err)
				return err
			}
		}

		err = b.SendMessage(c, "Отправь мне ссылку на репозиторий GitHub для отслеживания.")
		if err != nil {
			b.log.Error("TrackLink: Error sending message", "error", err)
			return err
		}

		err = b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingURL)
		if err != nil {
			b.log.Error("TrackLink: Error changing user state", "error", err)
			return err
		}

		return nil
	}
}
