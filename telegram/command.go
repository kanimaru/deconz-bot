package telegram

import (
	"telegram-deconz/deconz"
	"telegram-deconz/storage"
	"telegram-deconz/template"
)

type Command struct {
	bot            Bot
	deconzService  deconz.Service
	storageManager storage.Manager
	engine         template.Engine
}

func CreateCommand(bot Bot, deconzService deconz.Service, storageManager storage.Manager, engine template.Engine) *Command {
	return &Command{
		bot:            bot,
		deconzService:  deconzService,
		storageManager: storageManager,
		engine:         engine,
	}
}

func (c Command) GetStorage(message *Message) storage.Storage {
	return c.storageManager.Get(message.GetId())
}
