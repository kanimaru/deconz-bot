package telegram

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-deconz/deconz"
	"telegram-deconz/storage"
	"telegram-deconz/template"
)

type CommandFactory struct {
	bot            Bot
	deconzService  deconz.Service
	storageManager storage.Manager
	engine         template.Engine
}

func CreateCommandFactory(bot Bot, deconzService deconz.Service, storageManager storage.Manager, engine template.Engine) *CommandFactory {
	return &CommandFactory{
		bot:            bot,
		deconzService:  deconzService,
		storageManager: storageManager,
		engine:         engine,
	}
}

func (c CommandFactory) GetStorage(message *Message) storage.Storage {
	return c.storageManager.Get(message.GetId())
}

func (c CommandFactory) removeCommandMessage(message Message) {
	_, err := c.bot.Request(tgbotapi.NewDeleteMessage(message.GetChatId(), message.GetId()))
	if err != nil {
		log.Warningf("Can't delete the request")
	}
}

func (c CommandFactory) openInlineMessage(template string, data interface{}, message Message) {
	viewManager := CreateViewManager(c.bot, c.engine)
	_, newMessage := viewManager.Show(template, data, message)
	// Don't use storage from parameter it's for the command message.
	c.GetStorage(newMessage).Save("viewManager", viewManager)
}
