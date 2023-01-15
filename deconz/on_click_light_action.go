package deconz

import (
	"github.com/PerformLine/go-stockutil/log"
	"strconv"
	"strings"
	"telegram-deconz/bot"
	"telegram-deconz/storage"
	"telegram-deconz/template"
)

const ActionOn = "Action.On"
const ActionOff = "Action.Off"
const ActionColor = "Action.Color"
const ActionSetTemperature = "Action.SetTemperature"
const ActionSetBrightness = "Action.SetBrightness"

type LightActionOnClickHandler[Message bot.BaseMessage] struct {
	deconzService DeviceService

	selectedLights []string
	HandledActions []string
}

func CreateLightActionOnClickHandler[Message bot.BaseMessage](deconzService DeviceService) *LightActionOnClickHandler[Message] {
	return &LightActionOnClickHandler[Message]{
		deconzService: deconzService,
		HandledActions: []string{
			ActionOn,
			ActionOff,
			ActionColor,
			ActionSetTemperature,
			ActionSetBrightness,
		},
	}
}

func (l *LightActionOnClickHandler[Message]) CallAction(storage storage.Storage, message Message, target *template.Button) {
	views := storage.Get("viewManager").(bot.ViewManager[Message])
	lights := l.GetLights(target)
	l.selectedLights = lights
	if lights == nil {
		_ = views.Close(message)
		log.Errorf("No lights selected with this button use light:id or group:id")
	}

	switch *target.OnClick {
	case ActionOn:
		l.turnLight(lights, true)
		l.back(views, message)
	case ActionOff:
		l.turnLight(lights, false)
		l.back(views, message)
	case ActionColor:
		l.switchToColor(message)
	case ActionSetTemperature:
		l.setTemperature(target)
		l.back(views, message)
	case ActionSetBrightness:
		l.setBrightness(target)
		l.back(views, message)
	}

}

func (l *LightActionOnClickHandler[Message]) ReceiveMessage(message Message) {
	// TODO determine it for different chats otherwise one from a different chat can set the color too
	txt := message.GetText()
	if txt == "" {
		return
	}

	if len(txt) != 6 {
		return
	}

	l.deconzService.SetLightState(LightState{
		Color: txt,
	}, l.selectedLights...)
}

func (l *LightActionOnClickHandler[Message]) back(views bot.ViewManager[Message], message Message) {
	_, err := views.Back(message)
	if err != nil {
		_ = views.Close(message)
		log.Errorf("Can't back: %w", err)
	}
}

func (l *LightActionOnClickHandler[Message]) turnLight(lights []string, on bool) {
	l.deconzService.SetLightState(LightState{
		On: &on,
	}, lights...)
}

func (l *LightActionOnClickHandler[Message]) GetLights(button *template.Button) []string {
	for cur := &button.Element; cur != nil; cur = cur.Parent {
		data := cur.Data
		if strings.HasPrefix(data, "group:") {
			groupId := strings.Replace(data, "group:", "", 1)
			group := l.deconzService.GetGroup(groupId)
			return group.Lights
		} else if strings.HasPrefix(data, "light:") {
			light := strings.Replace(data, "light:", "", 1)
			return []string{light}
		}
	}
	return nil
}

func (l *LightActionOnClickHandler[Message]) switchToColor(message Message) {

}

func (l *LightActionOnClickHandler[Message]) setTemperature(target *template.Button) {
	temperature, _ := strconv.ParseInt(target.Data, 10, 8)
	if len(l.selectedLights) == 0 {
		return
	}

	lightId := l.selectedLights[0]
	light := l.deconzService.GetLight(lightId)
	colorTemp := int(float32(*light.Ctmin) + (float32(temperature)/100)*float32(*light.Ctmax-*light.Ctmin))

	l.deconzService.SetLightState(LightState{
		Temperature: &colorTemp,
	}, l.selectedLights...)
}

func (l *LightActionOnClickHandler[Message]) setBrightness(target *template.Button) {
	brightness, _ := strconv.ParseUint(target.Data, 10, 8)
	if len(l.selectedLights) == 0 {
		return
	}
	brightnessPtr := uint8(brightness)
	l.deconzService.SetLightState(LightState{
		Brightness: &brightnessPtr,
	}, l.selectedLights...)
}
