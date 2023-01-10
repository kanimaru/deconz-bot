package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kanimaru/godeconz/http"
)

type AllLightsCallback struct {
	LightSettingCallback
	content CallbackHandler
}

type AllLightData interface {
	GetGroup() string
}

type AllLightDataWrapper struct {
	allLightData  AllLightData
	deviceService DeviceService
}

func allLights(content CallbackHandler, data AllLightData, deviceService DeviceService, telegramService *TelegramService) CallbackHandler {
	newData := AllLightDataWrapper{
		allLightData:  data,
		deviceService: deviceService,
	}
	return &AllLightsCallback{
		LightSettingCallback: LightSettingCallback{
			telegramService: telegramService,
			data:            newData,
			deviceService:   deviceService,
		},
		content: content,
	}
}

func (a AllLightDataWrapper) GetLights() []string {
	lights := a.deviceService.GetLightsForGroup(a.GetGroup())
	lightIds := make([]string, 0, len(lights))
	for lightId := range lights {
		lightIds = append(lightIds, lightId)
	}
	return lightIds
}

func (a AllLightDataWrapper) GetGroup() string {
	return a.allLightData.GetGroup()
}

func (a *AllLightsCallback) initialize(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	msg, keyboard := a.content.initialize(update)

	deviceService := a.LightSettingCallback.deviceService
	groupName := a.LightSettingCallback.data.GetGroup()
	group := deviceService.GetGroup(groupName)
	lights := make([]http.LightResponseState, 0, len(group.Lights))
	for _, lightId := range group.Lights {
		light := deviceService.GetLight(lightId)
		lights = append(lights, light)
	}
	buttons, err := a.LightSettingCallback.getButtonsForLights(lights)
	if err != nil {
		return "No lights are reachable", tgbotapi.InlineKeyboardMarkup{}
	}

	var newKeyboard [][]tgbotapi.InlineKeyboardButton
	newKeyboard = append(newKeyboard, keyboard.InlineKeyboard...)
	newKeyboard = append(newKeyboard, tgbotapi.NewInlineKeyboardRow(buttons...))
	return msg, tgbotapi.NewInlineKeyboardMarkup(newKeyboard...)
}

func (a *AllLightsCallback) called(update tgbotapi.Update) CallbackHandler {
	callbackHandler := a.LightSettingCallback.called(update)
	if callbackHandler != nil {
		return callbackHandler
	}
	return a.content.called(update)
}
