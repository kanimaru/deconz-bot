package telegram

import (
	"telegram-deconz/bot"
	"telegram-deconz/storage"
	"telegram-deconz/view"
)

func (c Command) CreateLightCmd() bot.CommandDefinition[Message] {
	return bot.CommandDefinition[Message]{
		Description: "Let control lights and switches",
		Exec: func(storage storage.Storage, message Message) {
			c.removeCommandMessage(message)

			groupsMap := c.deconzService.GetGroups()
			groupsData := view.GroupsData{
				Groups: groupsMap,
			}

			c.openInlineMessage("groups.go.xml", groupsData, message)
		},
	}
}
