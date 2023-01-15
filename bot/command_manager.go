package bot

import (
	"github.com/PerformLine/go-stockutil/log"
	"telegram-deconz/storage"
)

type CommandDefinition[Message BaseMessage] struct {
	Description string
	Exec        func(storage storage.Storage, update Message)
}

type CommandManager[Message BaseMessage] struct {
	commands map[string]CommandDefinition[Message]
}

func CreateCommandManager[Message BaseMessage]() *CommandManager[Message] {
	return &CommandManager[Message]{
		commands: make(map[string]CommandDefinition[Message]),
	}
}

func (c *CommandManager[Message]) ReceiveMessage(storage storage.Storage, message Message) {
	command := message.GetCommand()
	if command == "" {
		return
	}

	commandDefinition, ok := c.commands[command]
	if !ok {
		log.Debugf("Command not found: %v", command)
		return
	}
	commandDefinition.Exec(storage, message)
}

func (c *CommandManager[Message]) AddCommand(command string, definition CommandDefinition[Message]) {
	c.commands[command] = definition
}

func (c *CommandManager[Message]) RemoveCommand(command string) CommandDefinition[Message] {
	definition := c.commands[command]
	delete(c.commands, command)
	return definition
}

func (c *CommandManager[Message]) GetCommands() map[string]CommandDefinition[Message] {
	return c.commands
}
