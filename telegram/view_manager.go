package telegram

import (
	"container/list"
	"encoding/xml"
	"errors"
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strings"
	"telegram-deconz/bot"
	"telegram-deconz/template"
)

type ViewManager struct {
	bot               Bot
	engine            template.Engine
	previousViewStack *list.List
	currentView       *template.View
	DebugView         bool
}

func CreateViewManager(bot Bot, engine template.Engine) bot.ViewManager[Message] {
	return &ViewManager{
		bot:               bot,
		engine:            engine,
		previousViewStack: list.New(),
		DebugView:         true,
	}
}

func (v *ViewManager) Open(view *template.View, message Message) (Message, error) {
	log.Infof("Open %w", view.Name)
	if v.DebugView {
		writeViewToFile(view, "temp/current_view.xml")
	}
	v.previousViewStack.PushBack(v.currentView)
	return v.changeKeyboard(view, message)
}

func (v *ViewManager) changeKeyboard(view *template.View, message Message) (Message, error) {
	v.currentView = view
	if message.GetText() != "" {
		msg := tgbotapi.NewMessage(message.GetChatId(), view.Text)
		msg.ReplyMarkup = GetInlineKeyboard(view)
		m, err := v.bot.Send(msg)

		return createMessageByMessage(&m), err
	} else if message.GetData() != "" {
		if strings.TrimSpace(view.Text) == "" {
			view.Text = view.Name
		}
		msg := tgbotapi.NewEditMessageTextAndMarkup(message.GetChatId(), message.GetId(), view.Text, GetInlineKeyboard(view))
		m, err := v.bot.Send(msg)
		return createMessageByMessage(&m), err
	}
	return Message{}, errors.New("message can't be handled")
}

func (v *ViewManager) Back(message Message) (*template.View, error) {
	latestElement := v.previousViewStack.Back()
	if latestElement == nil {
		log.Debugf("No history found - close view")
		return nil, v.Close(message)
	}
	latestView := latestElement.Value.(*template.View)
	log.Infof("Back to %w", latestView.Name)
	v.previousViewStack.Remove(latestElement)
	_, err := v.changeKeyboard(latestView, message)
	return latestView, err
}

func (v *ViewManager) Close(message Message) error {
	v.currentView = nil
	v.previousViewStack.Init()
	_, err := v.bot.Request(tgbotapi.NewDeleteMessage(message.GetChatId(), message.GetId()))
	return err
}

func (v *ViewManager) FindButton(id string) *template.Button {
	if v.currentView == nil {
		return nil
	}
	return v.currentView.FindButton(id)
}

func (v *ViewManager) Show(name string, data interface{}, message Message) (*template.View, *Message) {
	view, err := v.engine.Apply(name, data)
	if err != nil {
		_ = v.Close(message)
		log.Errorf("Problems with parsing template - %v: %w", name, err)
		return nil, nil
	}
	msg, err := v.Open(view, message)
	if err != nil {
		_ = v.Close(message)
		log.Errorf("Problems with displaying %v: %w", name, err)
	}
	return view, &msg
}

func GetInlineKeyboard(view *template.View) tgbotapi.InlineKeyboardMarkup {
	allButtons := make([][]tgbotapi.InlineKeyboardButton, 0, len(view.Row))
	for _, row := range view.Row {
		rowButtons := make([]tgbotapi.InlineKeyboardButton, 0, len(row.Button))
		for _, button := range row.Button {
			id := button.GetId()
			rowButtons = append(rowButtons, tgbotapi.NewInlineKeyboardButtonData(button.Label, id))
		}
		allButtons = append(allButtons, rowButtons)
	}
	return tgbotapi.NewInlineKeyboardMarkup(allButtons...)
}

func writeViewToFile(view *template.View, path string) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		log.Fatalf("Can't open file for writing view: %w", err)
	}
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "\t")
	err = encoder.Encode(view)
	if err != nil {
		log.Fatalf("Can't encode view to xml: %w", err)
	}
}
