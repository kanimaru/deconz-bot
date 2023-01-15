package deconz

import (
	"telegram-deconz/bot"
	"telegram-deconz/storage"
	"telegram-deconz/template"
	"telegram-deconz/view"
)

type LightsOnClickHandler[Message bot.BaseMessage] struct {
	deconzService DeviceService
}

func CreateLightsOnClickHandler[Message bot.BaseMessage](deconzService DeviceService) *LightsOnClickHandler[Message] {
	return &LightsOnClickHandler[Message]{
		deconzService: deconzService,
	}
}

func (l *LightsOnClickHandler[Message]) CallAction(storage storage.Storage, message Message, target *template.Button) {
	storage.Save("light", target.Data)
	groupId := storage.Get("group").(string)
	group := l.deconzService.GetGroup(groupId)
	views := storage.Get("viewManager").(bot.ViewManager[Message])

	light := l.deconzService.GetLight(target.Data)
	features := GetLightFeatures(light)

	lightData := view.LightData{
		GroupName:            group.Name,
		Id:                   target.Data,
		Name:                 light.Name,
		On:                   features.On,
		ColorAvailable:       features.HasColor,
		BrightnessAvailable:  features.HasBrightness,
		TemperatureAvailable: features.HasTemp,
	}

	views.Show("light.go.xml", lightData, message)
}
