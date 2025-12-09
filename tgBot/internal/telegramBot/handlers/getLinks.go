package handlers

import (
	"context"
	"fmt"
	"strings"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) AllLinks(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
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

		c.Send(msg)

		return nil
	}
}
