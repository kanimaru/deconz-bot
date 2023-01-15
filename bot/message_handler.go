package bot

import "container/list"

type BaseMessage interface {
	// GetId of the message
	GetId() int
	// GetChatId in which the message was sent
	GetChatId() int64
	// GetText of the message
	GetText() string
	// GetData of a message for example the ID of the button that got pressed in this message
	GetData() string
	// GetCommand without / or @bot name
	GetCommand() string
}

type MessageReceiver[Message BaseMessage] interface {
	// ReceiveMessage tries to find the clicked button and start handling it
	ReceiveMessage(message Message)
}

type OnDeleteMessageReceiver[Message BaseMessage] interface {
	// OnDeleteMessage called when a message got deleted so that according memory can be cleared too.
	OnDeleteMessage(message Message)
}

type MessageDistributor[Message BaseMessage] struct {
	messageReceiver  *list.List
	onDeleteReceiver *list.List
}

func CreateMessageDistributor[Message BaseMessage]() MessageDistributor[Message] {
	return MessageDistributor[Message]{
		messageReceiver:  list.New(),
		onDeleteReceiver: list.New(),
	}
}

func (b *MessageDistributor[Message]) AddMessageReceiver(receiver MessageReceiver[Message]) {
	b.messageReceiver.PushBack(receiver)
}

func (b *MessageDistributor[Message]) RemoveMessageReceiver(receiver MessageReceiver[Message]) {
	for el := b.messageReceiver.Front(); el != nil; el = el.Next() {
		if el.Value == receiver {
			b.messageReceiver.Remove(el)
			return
		}
	}
}

func (b *MessageDistributor[Message]) ReceiveMessage(message Message) {
	for el := b.messageReceiver.Front(); el != nil; el = el.Next() {
		el.Value.(MessageReceiver[Message]).ReceiveMessage(message)
	}
}

func (b *MessageDistributor[Message]) AddOnDeleteReceiver(receiver OnDeleteMessageReceiver[Message]) {
	b.onDeleteReceiver.PushBack(receiver)
}

func (b *MessageDistributor[Message]) RemoveOnDeleteReceiver(receiver OnDeleteMessageReceiver[Message]) {
	for el := b.onDeleteReceiver.Front(); el != nil; el = el.Next() {
		if el.Value == receiver {
			b.onDeleteReceiver.Remove(el)
			return
		}
	}
}

func (b *MessageDistributor[Message]) OnDeleteMessage(message Message) {
	for el := b.onDeleteReceiver.Front(); el != nil; el = el.Next() {
		el.Value.(OnDeleteMessageReceiver[Message]).OnDeleteMessage(message)
	}
}
