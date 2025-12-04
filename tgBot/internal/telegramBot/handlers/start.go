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
		btnHelp := menu.Text("Помощь")

		menu.Reply(menu.Row(btnTrack, btnUntrack), menu.Row(btnShowAll, btnHelp))

		err = b.useCase.CreateChat(ctx, c.Sender().ID)
		if err != nil {
			return c.Send("Вы уже зарегестрированы.", menu)
		}

		b.Bot.Send(c.Recipient(), "Добро пожаловать!\nВы успешно зарегистировались. Выберите команду:", menu)
		b.HelpMessage(c)
		return nil
	}
}
