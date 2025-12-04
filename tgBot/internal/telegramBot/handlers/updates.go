package handlers

import (
	"fmt"
	"linktracker/internal/domain"

	"gopkg.in/telebot.v4"
)

func (b *BotHandler) Updates(update *domain.UpdatedLink) error {
	msg := fmt.Sprintf("Произошло обновление по репозиторию %s\n\nОписание:\n%s", update.Link.URL, update.Link.Desc)

	for _, chatID := range update.ChatIDs {
		_, err := b.Bot.Send(telebot.ChatID(chatID), msg)
		if err != nil {
			return err
		}
	}

	return nil
}
