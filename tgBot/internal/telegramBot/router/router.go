package tgRouter

import (
	"context"
	"linktracker/internal/telegramBot/handlers"

	"gopkg.in/telebot.v4"
)

func Router(b *handlers.BotHandler, ctx context.Context) {
	b.Bot.Handle(telebot.OnText, b.MessageHandler(ctx))
	b.Bot.Handle("/start", b.Start(ctx))
	b.Bot.Handle("Привязать ссылку", b.TrackLink(ctx))
	b.Bot.Handle("Отвязать ссылку", b.UntrackLink(ctx))
	b.Bot.Handle("Список ссылок", b.AllLinks(ctx))
	b.Bot.Handle("Фильтр", b.GetFilteredLinks(ctx))
	b.Bot.Handle("Помощь", b.HelpMessage)
}
