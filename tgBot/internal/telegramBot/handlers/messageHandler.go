package handlers

import (
	"context"
	"linktracker/internal/domain"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) MessageHandler(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userInfo, err := b.useCase.GetUserState(ctx, c.Sender().ID)
		if err != nil {
			b.log.Error("MessageHandler: Error getting user state", "error", err)
			return err
		}
		switch userInfo.State {
		case domain.WaitingURl:
			_, isIt := b.LinkValidation(c.Text())
			if !isIt {
				b.Bot.Send(c.Recipient(), "Неправильно указана ссылка, попробуйте еще раз.")
				return nil
			}

			userInfo.URL = c.Text()

			b.log.Info("TrackLink: Waiting for link", "url", userInfo.URL)
			b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingDescription)
			b.Bot.Send(c.Recipient(), "Добавьте описание репозиторию\nЭтот этап необязательный, вы можете пропустить его написав 'skip'.")
		case domain.WaitingDescription:
			if c.Text() == "skip" {
				userInfo.Desc = ""
			} else {
				userInfo.Desc = c.Text()
			}
			b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)

			link := domain.Link{
				URL:    userInfo.URL,
				Desc:   userInfo.Desc,
				ChatID: c.Sender().ID,
			}

			err = b.useCase.AddLink(ctx, userInfo.UserID, link)
			if err != nil {
				b.Bot.Send(c.Recipient(), "Не удалось добавить ссылку.")
				b.log.Error("TrackLink: Error adding link", "error", err)
				return err
			}

			b.Bot.Send(c.Recipient(), "готово!")
		case domain.WaitingDelete:
			err := b.useCase.DeleteLink(ctx, c.Sender().ID, c.Text())
			if err != nil {
				b.log.Error("TrackLink: Error deleting link", "error", err)
				c.Send("Неправильно введено название ссылки, повторите еще раз")
				return err
			}

			b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
			b.Bot.Send(c.Recipient(), "Готово!\nВаша ссылка была успешно удалена.")
		case "/cancel":
			b.Cancel(ctx, userInfo, c)
		case "":
			b.Bot.Send(c.Recipient(), "Сперва нужно зарегистрироваться!\nВоспользуйтесь командой /start")
		}
		return nil
	}
}
