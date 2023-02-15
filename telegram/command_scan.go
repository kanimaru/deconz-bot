package telegram

import (
	"telegram-deconz/bot"
	"telegram-deconz/storage"
)

func (c CommandFactory) CreateScanCmd() bot.CommandDefinition[Message] {
	return bot.CommandDefinition[Message]{
		Description: "New device scanning",
		Exec: func(storage storage.Storage, message Message) {
			c.removeCommandMessage(message)
			c.openInlineMessage("scan.go.xml", nil, message)
		},
	}
}
