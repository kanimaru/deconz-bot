package bot

import (
	"github.com/PerformLine/go-stockutil/log"
	"telegram-deconz/storage"
	"telegram-deconz/template"
)

type OnClickHandler[Message BaseMessage] interface {
	CallAction(storage storage.Storage, message Message, target *template.Button)
}

type ActionManager[Message BaseMessage] struct {
	actions        map[string]OnClickHandler[Message]
	storageManager storage.Manager
}

func CreateBaseActionManager[Message BaseMessage](storageManager storage.Manager) *ActionManager[Message] {
	return &ActionManager[Message]{
		actions:        make(map[string]OnClickHandler[Message]),
		storageManager: storageManager,
	}
}

func (t *ActionManager[Message]) RegisterAction(handler OnClickHandler[Message], actions ...string) {
	for _, action := range actions {
		_, ok := t.actions[action]
		if ok {
			log.Warningf("Action '%v' got redefined!", action)
		}
		t.actions[action] = handler
	}
}

func (t *ActionManager[Message]) UnregisterAction(action string) OnClickHandler[Message] {
	handler, ok := t.actions[action]
	if ok {
		delete(t.actions, action)
		return handler
	}
	return nil
}

func (t *ActionManager[Message]) GetAction(data string) (OnClickHandler[Message], bool) {
	handler, ok := t.actions[data]
	return handler, ok
}

func (t *ActionManager[Message]) HandleAction(message Message, button *template.Button) {
	storage := t.storageManager.Get(message.GetId())
	if button.OnClick != nil {
		action, exists := t.GetAction(*button.OnClick)
		if exists {
			action.CallAction(storage, message, button)
		} else {
			log.Warningf("Action '%v' doesn't exists", *button.OnClick)
		}
	}
}
