package handlers

import (
	"context"
	"fmt"
	"linktracker/internal/domain"
	"strings"

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

			b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingTags)
			b.Bot.Send(c.Recipient(), "Добавьте теги для репозитория через запятую\nЭтот этап необязательный, вы можете пропустить его написав 'skip'.")
		case domain.WaitingTags:
			if c.Text() == "skip" {
				userInfo.Tags = ""
			} else {
				userInfo.Tags = c.Text()
			}

			link := domain.Link{
				URL:    userInfo.URL,
				Desc:   userInfo.Desc,
				Tags:   userInfo.Tags,
				ChatID: userInfo.UserID,
			}

			err = b.useCase.AddLink(ctx, userInfo.UserID, link)
			if err != nil {
				b.Bot.Send(c.Recipient(), "Не удалось добавить ссылку.")
				b.log.Error("TrackLink: Error adding link", "error", err)
				return err
			}

			b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
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
		case domain.WaitingFilter:
			links, err := b.useCase.GetFilteredLinks(ctx, c.Sender().ID, c.Text())
			if err != nil {
				b.log.Error("GetLinks: Error getting links", "error", err)
				return err
			}

			msg := "Ваши отслеживаемые ссылки с примененным фильтром: \n\n"

			if len(links) > 0 {
				for i := 0; i < len(links); i++ {
					urlParts := strings.Split(links[i].URL, "/")
					alias := urlParts[len(urlParts)-1]
					msg += fmt.Sprintf("%s: %s\nОписание репозитория: %s\nТеги репозитория: %s\n\n", alias, links[i].URL, links[i].Desc, links[i].Tags)
				}
			} else {
				msg = "Не удалось найти ссылки с данным тегом. Попробуйте еще раз."
			}

			c.Send(msg)
			b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
		case "/cancel":
			b.Cancel(ctx, userInfo, c)
		case "":
			b.Bot.Send(c.Recipient(), "Сперва нужно зарегистрироваться!\nВоспользуйтесь командой /start")
		}
		return nil
	}
}
