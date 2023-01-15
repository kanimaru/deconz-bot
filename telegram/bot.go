package telegram

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-deconz/bot"
)

type Bot struct {
	*tgbotapi.BotAPI
}

func CreateBot(apiKey string) Bot {
	api, err := tgbotapi.NewBotAPI(apiKey)
	if err != nil {
		log.Fatalf("Can't create telegram bot: %w", err)
	}
	return Bot{
		BotAPI: api,
	}
}

func (b Bot) HandleUpdates(receiver func(update Message)) {
	defer b.StopReceivingUpdates()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range b.GetUpdatesChan(u) {
		receiver(createMessageByUpdate(&update))
	}
}

func (b Bot) UpdateCommands(commandManager *bot.CommandManager[Message], scope tgbotapi.BotCommandScope) {
	commands := commandManager.GetCommands()
	var botCommands = make([]tgbotapi.BotCommand, 0, len(commands))
	for command, callback := range commands {
		botCommands = append(botCommands, tgbotapi.BotCommand{
			Command:     command,
			Description: callback.Description,
		})
	}

	_, err := b.Request(tgbotapi.NewSetMyCommandsWithScope(scope, botCommands...))
	if err != nil {
		log.Errorf("Can't update telegram bot commands: %v", err)
	}
}

type Message struct {
	id      int
	text    string
	data    string
	chatId  int64
	command string
}

func createMessageByMessage(message *tgbotapi.Message) Message {
	return Message{
		id:      message.MessageID,
		text:    message.Text,
		chatId:  message.Chat.ID,
		command: message.Command(),
	}
}

func createMessageByUpdate(update *tgbotapi.Update) Message {
	msg := Message{}

	if update.Message != nil {
		return createMessageByMessage(update.Message)
	}
	if update.CallbackQuery != nil {
		msg.id = update.CallbackQuery.Message.MessageID
		msg.data = update.CallbackQuery.Data
		msg.chatId = update.FromChat().ID
	}
	return msg
}

func (m Message) GetId() int {
	return m.id
}

func (m Message) GetText() string {
	return m.text
}

func (m Message) GetData() string {
	return m.data
}

func (m Message) GetChatId() int64 {
	return m.chatId
}

func (m Message) GetCommand() string {
	return m.command
}
