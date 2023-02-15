package main

import (
	"telegram-deconz/telegram"
)

// BaseMessageSender is a simple bridge to send messages to the default chat.
type BaseMessageSender struct {
	chatId int64
	bot    telegram.Bot
}

func CreateBaseMessageSender(chatId int64, bot telegram.Bot) *BaseMessageSender {
	return &BaseMessageSender{chatId: chatId, bot: bot}
}

func (b BaseMessageSender) SendMessage(message string) {
	b.bot.SendMessage(message, b.chatId)
}
