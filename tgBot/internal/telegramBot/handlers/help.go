package handlers

import "gopkg.in/telebot.v4"

func (b *BotHandler) HelpMessage(c telebot.Context) error {
	var parse telebot.ParseMode = "Markdown"
	b.Bot.Send(c.Recipient(), "Список действующий команд: "+
		"\n\n*Привязать ссылку* - привязать ссылку репозитория GiHub для отслеживание новый коммитов/проблем."+
		"\n\n*Отвязать ссылку* - отвязывает репозиторий от списка отслеживаемых ссылок."+
		"\n\n*Все ссылки* - отображает текущий список всех отслеживаемых репозиториев вами.", parse)
	return nil
}
