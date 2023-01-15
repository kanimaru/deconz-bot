package telegram

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-deconz/bot"
	"telegram-deconz/storage"
)

type ActionManager struct {
	*bot.ActionManager[Message]
	storageManager storage.Manager
	bot            *tgbotapi.BotAPI
}

func CreateActionManager(storageManager storage.Manager, api *tgbotapi.BotAPI) *ActionManager {
	return &ActionManager{
		ActionManager:  bot.CreateBaseActionManager[Message](storageManager),
		storageManager: storageManager,
		bot:            api,
	}
}

func (t *ActionManager) ReceiveMessage(message Message) {
	data := message.GetData()
	if message.GetData() == "" {
		return
	}
	storage := t.storageManager.Get(message.GetId())
	views := storage.Get("viewManager").(*ViewManager)

	button := views.FindButton(data)
	if button == nil {
		log.Errorf("Button not found.")
		return
	}

	t.HandleAction(message, button)
	if button.View != nil {
		_, err := views.Open(button.View, message)
		if err != nil {
			log.Errorf("There is a problem with changing the keyboards: %w", err)
		}
	}
}
