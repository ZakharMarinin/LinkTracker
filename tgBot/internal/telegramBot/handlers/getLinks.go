package handlers

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) AllLinks(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		_, err := b.useCase.GetUserState(ctx, c.Sender().ID)
		if err != nil {
			b.log.Warn("TrackLink: Error getting user state. Starting backoff...", "error", err)

			err = b.BackoffStart(ctx, c)
			if err != nil {
				b.log.Error("backoff start error", "error", err)
				return err
			}
		}

		links, err := b.useCase.GetLinks(ctx, c.Sender().ID)
		if err != nil {
			b.log.Error("GetLinks: Error getting links", "error", err)
			return err
		}

		msg := "Ваши отслеживаемые ссылки: \n\n"

		if len(links) > 0 {
			for i := 0; i < len(links); i++ {
				urlParts := strings.Split(links[i].URL, "/")
				alias := urlParts[len(urlParts)-1]
				msg += fmt.Sprintf("%s: %s\nОписание репозитория: %s\nТеги репозитория: %s\n\n", alias, links[i].URL, links[i].Desc, links[i].Tags)
			}
		} else {
			msg = "У вас пока что нет отслеживаемых ссылок."
		}

		err = b.SendMessage(c, msg)
		if err != nil {
			b.log.Error("Send: Error sending message", "error", err)
		}

		return nil
	}
}
