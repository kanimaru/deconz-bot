package template

import (
	"github.com/PerformLine/go-stockutil/log"
)

type MessageHandler[Message any] interface {
	// ReceiveMessage tries to find the clicked button and start handling it
	ReceiveMessage(message Message)

	// ChangeKeyboard will change the current keyboard for the user according to the view.
	// Message is the triggering message for that action.
	ChangeKeyboard(view *View, message Message) error

	// OnDeleteMessage called when a message got deleted so that according memory can be cleared too.
	OnDeleteMessage(message Message)
}

type OnClickHandler[Message any] interface {
	CallAction(storage Storage, message Message, target *Button)
}

type BaseMessageHandler[Message any] struct {
	buttons        map[string]*Button
	actions        map[string]OnClickHandler[Message]
	messageStorage StorageManager
}

func CreateBaseMessageHandler[Message any](storageManager StorageManager) *BaseMessageHandler[Message] {
	return &BaseMessageHandler[Message]{
		messageStorage: storageManager,
		actions:        make(map[string]OnClickHandler[Message]),
		buttons:        make(map[string]*Button),
	}
}

func (t *BaseMessageHandler[Message]) GetButton(data string) *Button {
	return t.buttons[data]
}

func (t *BaseMessageHandler[Message]) RegisterAction(action string, handler OnClickHandler[Message]) {
	_, ok := t.actions[action]
	if ok {
		log.Warningf("Action '%v' got redefined!", action)
	}
	t.actions[action] = handler
}

func (t *BaseMessageHandler[Message]) UnregisterAction(action string) OnClickHandler[Message] {
	handler, ok := t.actions[action]
	if ok {
		delete(t.actions, action)
		return handler
	}
	return nil
}

func (t *BaseMessageHandler[Message]) GetAction(data string) (OnClickHandler[Message], bool) {
	handler, ok := t.actions[data]
	return handler, ok
}

func (t *BaseMessageHandler[Message]) ReceiveMessage(message Message) {
	panic("implement me")
}

func (t *BaseMessageHandler[Message]) ChangeKeyboard(view *View, message Message) error {
	panic("implement me")
}

func (t *BaseMessageHandler[Message]) OnDeleteMessage(message Message) {
	panic("implement me")
}

func (t *BaseMessageHandler[Message]) onButtonClick(message Message, button *Button, storage Storage) {
	if button.OnClick != nil {
		action, exists := t.GetAction(*button.OnClick)
		if exists {
			action.CallAction(storage, message, button)
		} else {
			log.Warningf("Action '%v' doesn't exists", *button.OnClick)
		}
	}
	if button.View != nil {
		err := t.ChangeKeyboard(button.View, message)
		if err != nil {
			log.Errorf("There is a problem with changing the keyboards.")
		}
	}
}
