package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
)

type GroupCallback struct {
	message          string
	groups           map[string]string
	data             GroupData
	deviceService    DeviceService
	afterGroupSelect CallbackHandler
}

type GroupData interface {
	SetGroup(group string)
}

func createGroupContext(data GroupData, message string, deviceService DeviceService, afterGroupSelect CallbackHandler) CallbackHandler {
	return &GroupCallback{
		data:             data,
		message:          message,
		deviceService:    deviceService,
		afterGroupSelect: afterGroupSelect,
	}
}

func (g *GroupCallback) initialize(_ tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	g.groups = g.deviceService.GetGroups()
	groupButtonData := make([]buttonData, 0, len(g.groups))
	for data, group := range g.groups {
		groupButtonData = append(groupButtonData, buttonData{
			name: group,
			data: data,
		})
	}

	sort.Slice(groupButtonData, func(i, j int) bool {
		return groupButtonData[i].name < groupButtonData[j].name
	})

	buttonMatrix := createInlineButtonMatrix(groupButtonData, AutoColumn)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttonMatrix...)
	return g.message, keyboard
}

func (g *GroupCallback) messageReceived(_ tgbotapi.Update) {
	// ignored
}

func (g *GroupCallback) called(update tgbotapi.Update) CallbackHandler {
	g.data.SetGroup(update.CallbackQuery.Data)
	return g.afterGroupSelect
}

func (g *GroupCallback) cleanup(_ tgbotapi.Update) {
	// ignored
}
