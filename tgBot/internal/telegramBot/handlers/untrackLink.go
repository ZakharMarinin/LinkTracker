package handlers

import (
	"context"
	"fmt"
	"linktracker/internal/domain"
	"strings"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) UntrackLink(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userInfo, err := b.useCase.GetUserState(ctx, c.Sender().ID)
		if err != nil {
			b.log.Error("TrackLink: Error getting user state", "error", err)
			return err
		}

		if userInfo.State != domain.WaitingCommand {
			b.Bot.Send(c.Recipient(), "Сперва нужно зарегистрироваться!\nВоспользуйтесь командой /start")
			return nil
		}

		b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingDelete)

		msg := "Выберите ссылку для удаления, написав название репозитория: "

		links, err := b.useCase.GetLinks(ctx, c.Sender().ID)
		if err != nil {
			b.log.Error("GetLinks: Error getting links", "error", err)
			return err
		}

		for i := 0; i < len(links); i++ {
			urlParts := strings.Split(links[i].URL, "/")
			alias := urlParts[len(urlParts)-1]
			msg += fmt.Sprintf("%s: %s\n\n", alias, links[i].URL)
		}

		c.Send(msg)

		return nil
	}
}
