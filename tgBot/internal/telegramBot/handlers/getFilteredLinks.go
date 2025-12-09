package handlers

import (
	"context"
	"fmt"
	"linktracker/internal/domain"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) GetFilteredLinks(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userInfo, err := b.useCase.GetUserState(ctx, c.Sender().ID)
		fmt.Println(userInfo.State)
		if err != nil {
			b.log.Error("TrackLink: Error getting user state", "error", err)
			return err
		}

		if userInfo.State != domain.WaitingCommand {
			b.Bot.Send(c.Recipient(), "Сперва нужно зарегистрироваться!\nВоспользуйтесь командой /start")
			return nil
		}

		b.Bot.Send(c.Recipient(), "Напиши тег, по которому ты хочешь отфильтровать ссылки.")
		b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingFilter)

		return nil
	}
}
