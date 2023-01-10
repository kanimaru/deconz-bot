package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"sort"
)

type LightCallback struct {
	message          string
	data             LightData
	lights           map[string]string
	deviceService    DeviceService
	afterLightSelect CallbackHandler
}

type LightData interface {
	GetGroup() string
	SetLight(lightName string)
}

func createLightContext(data LightData, message string, deviceService DeviceService, afterLightSelect CallbackHandler) CallbackHandler {
	return &LightCallback{
		data:             data,
		message:          message,
		deviceService:    deviceService,
		afterLightSelect: afterLightSelect,
	}
}

func (l *LightCallback) initialize(_ tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	group := l.data.GetGroup()
	l.lights = l.deviceService.GetLightsForGroup(group)

	lightButtonData := make([]buttonData, 0, len(l.lights))
	for data, light := range l.lights {
		lightButtonData = append(lightButtonData, buttonData{
			name: light,
			data: data,
		})
	}

	sort.Slice(lightButtonData, func(i, j int) bool {
		return lightButtonData[i].name < lightButtonData[j].name
	})

	buttonMatrix := createInlineButtonMatrix(lightButtonData, AutoColumn)
	keyboard := tgbotapi.NewInlineKeyboardMarkup(buttonMatrix...)
	return l.message, keyboard
}

func (l *LightCallback) messageReceived(_ tgbotapi.Update) {
}

func (l *LightCallback) called(update tgbotapi.Update) CallbackHandler {
	l.data.SetLight(update.CallbackQuery.Data)
	return l.afterLightSelect
}

func (l *LightCallback) cleanup(_ tgbotapi.Update) {
}
