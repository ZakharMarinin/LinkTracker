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

		if c.Message().Text == "/cancel" {
			err = b.Cancel(ctx, userInfo, c)
			if err != nil {
				b.log.Error("MessageHandler: Error sending message", "error", err)
				return err
			}
		}

		switch userInfo.State {
		case domain.WaitingURL:
			err = b.WaitingURL(ctx, c, userInfo)
			if err != nil {
				b.log.Error("MessageHandler: Error getting waiting URL", "error", err)
				return err
			}

		case domain.WaitingDescription:
			err = b.WaitingDescription(ctx, c, userInfo)
			if err != nil {
				b.log.Error("MessageHandler: Error getting waiting description", "error", err)
				return err
			}

		case domain.WaitingTags:
			err = b.WaitingTags(ctx, c, userInfo)
			if err != nil {
				b.log.Error("MessageHandler: Error getting waiting tags", "error", err)
				return err
			}

		case domain.WaitingDelete:
			err = b.WaitingDelete(ctx, c, userInfo)
			if err != nil {
				b.log.Error("MessageHandler: Error getting waiting delete", "error", err)
				return err
			}

		case domain.WaitingFilter:
			err = b.WaitingFilter(ctx, c, userInfo)
			if err != nil {
				b.log.Error("MessageHandler: Error getting waiting filter", "error", err)
				return err
			}
		}

		return nil
	}
}

func (b *BotHandler) WaitingURL(ctx context.Context, c telebot.Context, userInfo *domain.UserStateInfo) error {
	_, isIt := b.LinkValidation(c.Text())
	if !isIt {
		err := b.SendMessage(c, "Неправильно указана ссылка, попробуйте еще раз.")
		if err != nil {
			b.log.Error("MessageHandler: Error sending message", "error", err)
			return err
		}

		return nil
	}

	userInfo.URL = c.Text()

	b.log.Info("TrackLink: Waiting for link", "url", userInfo.URL)
	err := b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingDescription)

	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	err = b.SendMessage(c, "Добавьте описание репозиторию\nЭтот этап необязательный, вы можете пропустить его написав 'skip'.")
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	return nil
}

func (b *BotHandler) WaitingDescription(ctx context.Context, c telebot.Context, userInfo *domain.UserStateInfo) error {
	if c.Text() == "skip" {
		userInfo.Desc = ""
	} else {
		userInfo.Desc = c.Text()
	}

	err := b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingTags)
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	err = b.SendMessage(c, "Добавьте теги для репозитория через запятую\n"+
		"Этот этап необязательный, вы можете пропустить его написав 'skip'.")
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	return nil
}

func (b *BotHandler) WaitingTags(ctx context.Context, c telebot.Context, userInfo *domain.UserStateInfo) error {
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

	err := b.useCase.AddLink(ctx, userInfo.UserID, &link)
	if err != nil {
		err = b.SendMessage(c, "Не удалось добавить ссылку.")
		if err != nil {
			b.log.Error("MessageHandler: Error sending message", "error", err)
			return err
		}

		b.log.Error("TrackLink: Error adding link", "error", err)

		return err
	}

	err = b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	err = b.SendMessage(c, "Готово!")
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	return nil
}

func (b *BotHandler) WaitingDelete(ctx context.Context, c telebot.Context, userInfo *domain.UserStateInfo) error {
	err := b.useCase.DeleteLink(ctx, c.Sender().ID, c.Text())
	if err != nil {
		b.log.Error("TrackLink: Error deleting link", "error", err)
		err = b.SendMessage(c, "Неправильно введено название ссылки, повторите еще раз")

		if err != nil {
			b.log.Error("TrackLink: Error sending message", "error", err)
			return err
		}

		return err
	}

	err = b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	err = b.SendMessage(c, "Готово!\nВаша ссылка была успешно удалена.")
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	return nil
}

func (b *BotHandler) WaitingFilter(ctx context.Context, c telebot.Context, userInfo *domain.UserStateInfo) error {
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

	err = b.SendMessage(c, msg)
	if err != nil {
		b.log.Error("Send: Error sending message", "error", err)
		return err
	}

	err = b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
	if err != nil {
		b.log.Error("MessageHandler: Error sending message", "error", err)
		return err
	}

	return nil
}
