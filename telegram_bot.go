package main

import (
	"fmt"
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type CommandCallback struct {
	description string
	exec        func(update tgbotapi.Update)
}

type TelegramService struct {
	*tgbotapi.BotAPI
	callbackHandler CallbackHandler
	logger          Logger
	commands        map[string]CommandCallback
	chatId          int64
}

type CallbackHandler interface {
	// initialize the keyboard layout
	initialize(tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup)
	// normal message got received and can be handled
	messageReceived(tgbotapi.Update)
	// called when callback hits, can optional return a new handler
	called(tgbotapi.Update) CallbackHandler
	// cleanup the old data if necessary
	cleanup(tgbotapi.Update)
}

func createTelegramClient() *TelegramService {
	apiKey := getEnv("TELEGRAM_API_KEY", "")
	defaultChatId := getEnv("TELEGRAM_CHAT_ID", "")
	chatId, err := strconv.ParseInt(defaultChatId, 10, 64)
	if err != nil {
		log.Fatalf("Can't parse TELEGRAM_CHAT_ID to int")
	}

	bot, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Panicf("Telegram can't be initialized... %v", err)
	}

	log.Noticef("Telegram Bot initialized...")
	return &TelegramService{
		BotAPI:          bot,
		callbackHandler: nil,
		chatId:          chatId,
		logger:          Logger{},
		commands:        make(map[string]CommandCallback),
	}

}

func (t *TelegramService) handleUpdates(close chan bool) {
	defer t.StopReceivingUpdates()

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range t.GetUpdatesChan(u) {
		select {
		case <-close:
			return
		default:
			if update.Message != nil {
				if update.Message.IsCommand() {
					t.handleCommands(update)
				} else {
					t.callbackHandler.messageReceived(update)
				}
			} else if update.CallbackQuery != nil && t.callbackHandler != nil {
				newHandler := t.callbackHandler.called(update)
				if newHandler != nil {
					t.changeContext(update, newHandler)
				}
			}
		}
	}
}

func (t *TelegramService) changeContext(update tgbotapi.Update, callbackHandler CallbackHandler) {
	if t.callbackHandler != nil {
		t.callbackHandler.cleanup(update)
	}
	t.callbackHandler = callbackHandler
	msg, keyboard := callbackHandler.initialize(update)

	if update.CallbackQuery != nil {
		t.SetInlineKeyboard(update, msg, keyboard)
	} else {
		t.SetKeyboard(update, msg, keyboard)
	}
}

func (t *TelegramService) AddCommand(command string, callback CommandCallback) {
	t.commands[command] = callback
}

func (t *TelegramService) handleCommands(update tgbotapi.Update) {
	cmd := update.Message.Command()
	callback, ok := t.commands[cmd]
	if ok {
		callback.exec(update)
	}
	t.logger.Debugf("Command not found: %v", cmd)
}

func (t *TelegramService) cleanupInlineError(update tgbotapi.Update, err error) {
	var msgId int
	if update.CallbackQuery != nil {
		msgId = update.CallbackQuery.Message.MessageID
	} else {
		msgId = update.Message.MessageID
	}
	deleteMessage := tgbotapi.NewDeleteMessage(update.FromChat().ID, msgId)
	t.SendIgnoreError(deleteMessage)
	text := fmt.Sprintf("Sorry command isn't working right now try it later. \n Cause: %v", err)
	message := tgbotapi.NewMessage(update.FromChat().ID, text)
	t.SendIgnoreError(message)
	t.callbackHandler.cleanup(update)
}

func (t *TelegramService) SendIgnoreError(chattable tgbotapi.Chattable) {
	_, err := t.Send(chattable)
	if err != nil {
		t.logger.Errorf("Can't send message: %v", err)
	}
}

func (t *TelegramService) SetKeyboard(update tgbotapi.Update, message string, keyboard tgbotapi.InlineKeyboardMarkup) {
	delMsg := tgbotapi.NewDeleteMessage(update.Message.Chat.ID, update.Message.MessageID)
	_, err := t.Request(delMsg)
	if err != nil {
		t.logger.Errorf("Can't send message: %v", err)
	}
	m := tgbotapi.NewMessage(update.Message.Chat.ID, message)
	m.ReplyMarkup = keyboard
	_, err = t.Send(m)
	if err != nil {
		t.cleanupInlineError(update, err)
	}
}

func (t *TelegramService) SetInlineKeyboard(update tgbotapi.Update, message string, keyboard tgbotapi.InlineKeyboardMarkup) {
	msg := update.CallbackQuery.Message
	newMessage := tgbotapi.NewEditMessageTextAndMarkup(msg.Chat.ID, msg.MessageID, message, keyboard)
	_, err := t.Send(newMessage)
	if err != nil {
		log.Errorf("Can't set inline keyboard: %v", err)
	}
}

func (t *TelegramService) UpdateCommands() {
	var botCommands = make([]tgbotapi.BotCommand, 0, len(t.commands))
	for command, callback := range t.commands {
		botCommands = append(botCommands, tgbotapi.BotCommand{
			Command:     command,
			Description: callback.description,
		})
	}

	scopeChat := tgbotapi.NewBotCommandScopeChat(t.chatId)
	_, err := t.Request(tgbotapi.NewSetMyCommandsWithScope(scopeChat, botCommands...))
	if err != nil {
		t.logger.Errorf("Can't update telegram bot commands: %v", err)
	} else {
		log.Notice("Commands are up to date.")
	}
}

// AutoColumn can be used in createInlineButtonMatrix to auto create 2 or 3 columns depending on the data to render
const AutoColumn = -1

type buttonData struct {
	name string
	data string
}

func createInlineButtonMatrix(buttons []buttonData, columns int) [][]tgbotapi.InlineKeyboardButton {
	if columns == AutoColumn {
		groupLen := len(buttons)
		columns = 1
		if groupLen%2 == 0 {
			columns = 2
		} else if groupLen%3 == 0 {
			columns = 3
		}
	}
	buttonColumn := make([]tgbotapi.InlineKeyboardButton, columns)
	buttonRows := make([][]tgbotapi.InlineKeyboardButton, len(buttons)/columns)
	colIndex := 0
	rowIndex := 0

	for _, button := range buttons {
		buttonColumn[colIndex] = tgbotapi.NewInlineKeyboardButtonData(button.name, button.data)
		colIndex++
		colIndex %= columns
		if colIndex == 0 {
			buttonRows[rowIndex] = buttonColumn
			rowIndex++
			buttonColumn = make([]tgbotapi.InlineKeyboardButton, columns)
		}
	}
	return buttonRows
}

func make100Keyboard() [][]tgbotapi.InlineKeyboardButton {
	values := make([]buttonData, 0, 5)
	for val := 0; val <= 100; val += 25 {
		values = append(values, buttonData{
			name: fmt.Sprintf("%d %%", val),
			data: fmt.Sprintf("%d", val),
		})
	}
	buttons := createInlineButtonMatrix(values, 5)
	return buttons
}
