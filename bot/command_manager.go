package bot

import (
	"github.com/PerformLine/go-stockutil/log"
)

type CommandDefinition[Message BaseMessage] struct {
	Description string
	Exec        func(update Message)
}

type CommandManager[Message BaseMessage] struct {
	commands map[string]CommandDefinition[Message]
}

func CreateCommandManager[Message BaseMessage]() *CommandManager[Message] {
	return &CommandManager[Message]{
		commands: make(map[string]CommandDefinition[Message]),
	}
}

func (c *CommandManager[Message]) ReceiveMessage(message Message) {
	command := message.GetCommand()
	if command != "" {
		c.handleCommand(message)
	}
}

func (c *CommandManager[Message]) handleCommand(message Message) {
	log.Infof("Msg ID: %v", message.GetId())
	command := message.GetCommand()
	commandDefinition, ok := c.commands[command]
	if !ok {
		log.Debugf("Command not found: %v", command)
		return
	}
	commandDefinition.Exec(message)
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
