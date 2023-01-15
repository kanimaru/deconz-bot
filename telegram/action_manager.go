package telegram

import (
	"github.com/PerformLine/go-stockutil/log"
	"telegram-deconz/bot"
	"telegram-deconz/storage"
)

type ActionManager struct {
	*bot.ActionManager[Message]
	storageManager storage.Manager
}

func CreateActionManager(storageManager storage.Manager) *ActionManager {
	return &ActionManager{
		ActionManager:  bot.CreateBaseActionManager[Message](storageManager),
		storageManager: storageManager,
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
