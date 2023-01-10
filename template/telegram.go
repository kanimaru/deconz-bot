package template

import (
	"errors"
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramTemplate struct {
	BaseMessageHandler[*tgbotapi.Update]
	bot *tgbotapi.BotAPI
}

func CreateTelegramTemplate(bot *tgbotapi.BotAPI) *TelegramTemplate {
	return &TelegramTemplate{
		BaseMessageHandler: *CreateBaseMessageHandler[*tgbotapi.Update](CreateInMemoryStorage()),
		bot:                bot,
	}
}

func (t *TelegramTemplate) GetInlineKeyboard(view *View) tgbotapi.InlineKeyboardMarkup {
	allButtons := make([][]tgbotapi.InlineKeyboardButton, 0, len(view.Row))
	for _, row := range view.Row {
		rowButtons := make([]tgbotapi.InlineKeyboardButton, 0, len(row.Button))
		for _, button := range row.Button {
			id := button.getId()
			rowButtons = append(rowButtons, tgbotapi.NewInlineKeyboardButtonData(button.Label, id))
			t.buttons[id] = &button
		}
		allButtons = append(allButtons, rowButtons)
	}
	return tgbotapi.NewInlineKeyboardMarkup(allButtons...)
}

func (t *TelegramTemplate) OnDeleteMessage(message *tgbotapi.Message) {
	t.messageStorage.Remove(message.MessageID)
	log.Debugf("Removed storage for: %v", message.MessageID)
}

func (t *TelegramTemplate) ReceiveMessage(update *tgbotapi.Update) {
	storage := t.messageStorage.Get(getMessageId(update))
	data := update.CallbackData()
	if data != "" {
		button := t.GetButton(data)
		t.onButtonClick(update, button, storage)
	}
}

func (t *TelegramTemplate) ChangeKeyboard(view *View, update *tgbotapi.Update) error {
	if update.Message != nil {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, view.Text)
		msg.ReplyMarkup = t.GetInlineKeyboard(view)
		_, err := t.bot.Send(msg)
		return err
	} else if update.CallbackQuery != nil {
		message := update.CallbackQuery.Message
		msg := tgbotapi.NewEditMessageTextAndMarkup(message.Chat.ID, message.MessageID, view.Text, t.GetInlineKeyboard(view))
		_, err := t.bot.Send(msg)
		return err
	}
	return errors.New("message can't be handled")
}

type BaseTelegramOnClickHandler struct {
	//telegramTemplate *TelegramTemplate
}

func getMessageId(update *tgbotapi.Update) int {
	var messageId int
	if update.Message != nil {
		messageId = update.Message.MessageID
	}
	if update.CallbackQuery.Message != nil {
		messageId = update.CallbackQuery.Message.MessageID
	}
	return messageId
}

/**
1. Action handler gets message and cares about the context
*/
