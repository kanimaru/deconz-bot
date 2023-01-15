package telegram

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

func (c Command) removeCommandMessage(message Message) {
	_, err := c.bot.Request(tgbotapi.NewDeleteMessage(message.GetChatId(), message.GetId()))
	if err != nil {
		log.Warningf("Can't delete the request")
	}
}

func (c Command) openInlineMessage(template string, data interface{}, message Message) {
	viewManager := CreateViewManager(c.bot, c.engine)
	_, newMessage := viewManager.Show(template, data, message)
	// Don't use storage from parameter it's for the command message.
	c.GetStorage(newMessage).Save("viewManager", viewManager)
}
