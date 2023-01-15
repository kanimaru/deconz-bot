package telegram

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"telegram-deconz/bot"
	"telegram-deconz/view"
)

func (c Command) CreateLightCmd() bot.CommandDefinition[Message] {
	return bot.CommandDefinition[Message]{
		Description: "Let control lights and switches",
		Exec: func(message Message) {
			_, err := c.bot.Request(tgbotapi.NewDeleteMessage(message.GetChatId(), message.GetId()))
			if err != nil {
				log.Warningf("Can't delete the request")
			}

			groupsMap := c.deviceService.GetGroups()
			groupsData := view.GroupsData{
				Groups: groupsMap,
			}

			viewManager := CreateViewManager(c.bot, c.engine)
			_, newMessage := viewManager.Show("groups.go.xml", groupsData, message)
			c.GetStorage(newMessage).Save("viewManager", viewManager)
		},
	}
}
