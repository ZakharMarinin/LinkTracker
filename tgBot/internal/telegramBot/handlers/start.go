package handlers

import (
	"context"
	"linktracker/internal/domain"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) Start(ctx context.Context) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		userInfo := &domain.UserStateInfo{
			UserID: c.Sender().ID,
		}

		err := b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
		if err != nil {
			return err
		}

		menu := &telebot.ReplyMarkup{ResizeKeyboard: true}
		btnTrack := menu.Text("Привязать ссылку")
		btnUntrack := menu.Text("Отвязать ссылку")
		btnShowAll := menu.Text("Список ссылок")
		btnFilter := menu.Text("Фильтр")
		btnHelp := menu.Text("Помощь")

		menu.Reply(menu.Row(btnTrack, btnUntrack), menu.Row(btnShowAll, btnFilter), menu.Row(btnHelp))

		err = b.useCase.CreateChat(ctx, c.Sender().ID)
		if err != nil {
			return c.Send("Вы уже зарегистрированы.", menu)
		}

		err = b.SendMessage(c, "Добро пожаловать!\nВы успешно зарегистрировались. Выберите команду:", menu)
		if err != nil {
			b.log.Error("Send: Error sending message", "error", err)
			return err
		}

		err = b.HelpMessage(c)
		if err != nil {
			b.log.Error("HelpMessage: Error sending message", "error", err)
			return err
		}

		return nil
	}
}

func (b *BotHandler) BackoffStart(ctx context.Context, c telebot.Context) error {
	userInfo := &domain.UserStateInfo{
		UserID: c.Sender().ID,
	}

	err := b.useCase.ChangeUserState(ctx, userInfo, domain.WaitingCommand)
	if err != nil {
		return err
	}

	err = b.useCase.CreateChat(ctx, c.Sender().ID)
	if err != nil {
		return b.SendMessage(c, "Вы уже зарегистрировались.")
	}

	return nil
}
