package bot

type Bot[Message BaseMessage] interface {
	SendMessage(chatId int, text string) (Message, error)
}

type BaseBot struct {
}
