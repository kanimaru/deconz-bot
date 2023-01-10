package main

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ViewCallback struct {
	telegramService  *TelegramService
	previousCallback CallbackHandler
	content          CallbackHandler
}

func asView(content CallbackHandler, telegramService *TelegramService) CallbackHandler {
	return &ViewCallback{
		telegramService: telegramService,
		content:         content,
	}
}

func (v *ViewCallback) initialize(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	msg, keyboard := v.content.initialize(update)

	var newKeyboard [][]tgbotapi.InlineKeyboardButton
	closeButton := tgbotapi.NewInlineKeyboardButtonData("Close", "close")
	var controlButtons []tgbotapi.InlineKeyboardButton
	if v.previousCallback != nil {
		backButton := tgbotapi.NewInlineKeyboardButtonData("Back", "back")
		controlButtons = tgbotapi.NewInlineKeyboardRow(closeButton, backButton)
	} else {
		controlButtons = tgbotapi.NewInlineKeyboardRow(closeButton)
	}

	newKeyboard = append(newKeyboard, controlButtons)
	newKeyboard = append(newKeyboard, keyboard.InlineKeyboard...)
	return msg, tgbotapi.NewInlineKeyboardMarkup(newKeyboard...)
}

func (v *ViewCallback) messageReceived(update tgbotapi.Update) {
	v.content.messageReceived(update)
}

func (v *ViewCallback) called(update tgbotapi.Update) CallbackHandler {
	switch update.CallbackQuery.Data {
	case "close":
		v.Close(update)
		return nil
	case "back":
		return v.previousCallback
	}
	newHandler := v.content.called(update)
	// Don't create an endless loop on the way back
	if newHandler == v.previousCallback || newHandler == v {
		return newHandler
	}

	newViewCallback, ok := newHandler.(*ViewCallback)
	if ok {
		newViewCallback.previousCallback = v
	}
	return newHandler
}

func (v *ViewCallback) cleanup(update tgbotapi.Update) {
	v.content.cleanup(update)
}

func (v *ViewCallback) Close(update tgbotapi.Update) {
	deleteMessage := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
	_, err := v.telegramService.Request(deleteMessage)
	if err != nil {
		log.Fatalf("Can't send delete message: %v", err)
	}
}
