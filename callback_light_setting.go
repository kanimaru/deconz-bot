package main

import (
	"errors"
	"fmt"
	"github.com/PerformLine/go-stockutil/log"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kanimaru/godeconz/http"
	"github.com/lucasb-eyer/go-colorful"
	"strconv"
)

type LightSettingCallback struct {
	telegramService *TelegramService
	data            LightSettingData
	deviceService   DeviceService
}

type LightSettingData interface {
	GetLights() []string
	GetGroup() string
}

func createLightSettingCallback(data LightSettingData, deviceService DeviceService, telegramService *TelegramService) *LightSettingCallback {
	return &LightSettingCallback{
		data:            data,
		deviceService:   deviceService,
		telegramService: telegramService,
	}
}

func (l *LightSettingCallback) getButtonsForLights(lights []http.LightResponseState) ([]tgbotapi.InlineKeyboardButton, error) {
	isOn := false
	hasBrightness := false
	hasColor := false
	hasTemp := false
	reachable := false
	for _, light := range lights {
		reachable = reachable || (light.State.Reachable != nil && *light.State.Reachable)
		hasTemp = hasTemp || (light.Ctmin != nil && light.Ctmax != nil && (*light.Ctmin != *light.Ctmax))
		hasColor = hasColor || (light.Hascolor != nil && *light.Hascolor)
		isOn = isOn || (light.State.On != nil && *light.State.On)
		hasBrightness = hasBrightness || light.State.Bri != nil
	}

	if !reachable {
		return nil, errors.New("no light is reachable")
	}

	var buttons []tgbotapi.InlineKeyboardButton

	if isOn {
		offBtn := tgbotapi.NewInlineKeyboardButtonData("Off", "off")
		buttons = append(buttons, offBtn)
	} else {
		onBtn := tgbotapi.NewInlineKeyboardButtonData("On", "on")
		buttons = append(buttons, onBtn)
	}

	if hasColor {
		colorBtn := tgbotapi.NewInlineKeyboardButtonData("Color", "color")
		buttons = append(buttons, colorBtn)
	}

	if hasTemp {
		tempBtn := tgbotapi.NewInlineKeyboardButtonData("Temperature", "temp")
		buttons = append(buttons, tempBtn)
	}

	if hasBrightness {
		brightnessBtn := tgbotapi.NewInlineKeyboardButtonData("Brightness", "bright")
		buttons = append(buttons, brightnessBtn)
	}
	return buttons, nil
}

func (l *LightSettingCallback) initialize(tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	lights := make([]http.LightResponseState, 0, len(l.data.GetLights()))
	for _, lightId := range l.data.GetLights() {
		light := l.deviceService.GetLight(lightId)
		lights = append(lights, light)
	}
	buttons, err := l.getButtonsForLights(lights)
	if err != nil {
		return "Lights are not reachable.", tgbotapi.InlineKeyboardMarkup{}
	}

	group := l.deviceService.GetGroup(l.data.GetGroup())
	txt := fmt.Sprintf("Lights of %v", group.Name)
	if len(lights) == 1 {
		txt = fmt.Sprintf("Light %v / %v", group.Name, lights[0].Name)
	}
	return txt, tgbotapi.NewInlineKeyboardMarkup(buttons)
}

func (l *LightSettingCallback) messageReceived(tgbotapi.Update) {
}

var off = false
var on = true

func (l *LightSettingCallback) called(update tgbotapi.Update) CallbackHandler {
	switch update.CallbackData() {
	case "off":
		l.deviceService.SetLightState(LightState{
			on: &off,
		}, l.data.GetLights()...)
		return l.telegramService.callbackHandler
	case "on":
		l.deviceService.SetLightState(LightState{
			on: &on,
		}, l.data.GetLights()...)
		return l.telegramService.callbackHandler
	case "bright":
		return asView(&BrightnessCallback{
			data:                    l.data,
			deviceService:           l.deviceService,
			previousCallbackHandler: l.telegramService.callbackHandler,
		}, l.telegramService)
	case "temp":
		return asView(&TemperatureCallback{
			data:                    l.data,
			deviceService:           l.deviceService,
			previousCallbackHandler: l.telegramService.callbackHandler,
		}, l.telegramService)
	case "color":
		return asView(&ColorCallback{
			telegramService:         l.telegramService,
			data:                    l.data,
			deviceService:           l.deviceService,
			previousCallbackHandler: l.telegramService.callbackHandler,
			currentMessage:          update,
		}, l.telegramService)
	}
	return nil
}

func (l *LightSettingCallback) cleanup(tgbotapi.Update) {}

//
// Brightness
//

type BrightnessCallback struct {
	data                    LightSettingData
	deviceService           DeviceService
	previousCallbackHandler CallbackHandler
}

func (b *BrightnessCallback) initialize(tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	return "Set Brightness", tgbotapi.NewInlineKeyboardMarkup(make100Keyboard()...)
}

func (b *BrightnessCallback) messageReceived(tgbotapi.Update) {}

func (b *BrightnessCallback) called(update tgbotapi.Update) CallbackHandler {
	brightness, err := strconv.ParseInt(update.CallbackData(), 10, 8)
	brightnessPtr := uint8(brightness)
	onState := false
	if brightness > 0 {
		onState = true
	}
	if err != nil {
		log.Fatalf("Can't parse brightness level %v: %v", update.CallbackData(), err)
	}
	b.deviceService.SetLightState(LightState{
		on:         &onState,
		brightness: &brightnessPtr,
	}, b.data.GetLights()...)
	return b.previousCallbackHandler
}

func (b *BrightnessCallback) cleanup(tgbotapi.Update) {}

//
// Temp
//

type TemperatureCallback struct {
	data                    LightSettingData
	deviceService           DeviceService
	previousCallbackHandler CallbackHandler
}

func (t *TemperatureCallback) initialize(tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	return "Set Color Temperature (0% coldest / 100% warmest)", tgbotapi.NewInlineKeyboardMarkup(make100Keyboard()...)
}

func (t *TemperatureCallback) messageReceived(tgbotapi.Update) {}

func (t *TemperatureCallback) called(update tgbotapi.Update) CallbackHandler {
	for _, l := range t.data.GetLights() {
		light := t.deviceService.GetLight(l)

		temperature, err := strconv.ParseInt(update.CallbackData(), 10, 8)
		temperaturePtr := int(float32(*light.Ctmin) + (float32(temperature)/100)*float32(*light.Ctmax-*light.Ctmin))

		if err != nil {
			log.Fatalf("Can't parse brightness level %v: %v", update.CallbackData(), err)
		}
		t.deviceService.SetLightState(LightState{
			temperature: &temperaturePtr,
		}, l)
	}
	return t.previousCallbackHandler
}

func (t *TemperatureCallback) cleanup(tgbotapi.Update) {}

//
// Color
//

type ColorCallback struct {
	telegramService         *TelegramService
	data                    LightSettingData
	deviceService           DeviceService
	previousCallbackHandler CallbackHandler
	currentMessage          tgbotapi.Update
}

func (c *ColorCallback) initialize(update tgbotapi.Update) (string, tgbotapi.InlineKeyboardMarkup) {
	c.currentMessage = update
	return "Please enter a RGB Color (ex: #993333):", tgbotapi.InlineKeyboardMarkup{}
}

func (c *ColorCallback) messageReceived(update tgbotapi.Update) {
	color := update.Message.Text
	_, err := colorful.Hex(color)
	if err != nil {
		log.Debugf("Not the right message")
		return
	}

	c.deviceService.SetLightState(LightState{color: color}, c.data.GetLights()...)
	c.telegramService.changeContext(c.currentMessage, c.previousCallbackHandler)
}

func (c *ColorCallback) called(tgbotapi.Update) CallbackHandler {
	return nil
}

func (c *ColorCallback) cleanup(tgbotapi.Update) {
}
