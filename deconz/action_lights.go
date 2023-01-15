package deconz

import (
	"telegram-deconz/bot"
	"telegram-deconz/storage"
	"telegram-deconz/template"
	"telegram-deconz/view"
)

type LightsAction[Message bot.BaseMessage] struct {
	deconzService Service
}

func CreateLightsAction[Message bot.BaseMessage](deconzService Service) *LightsAction[Message] {
	return &LightsAction[Message]{
		deconzService: deconzService,
	}
}

func (l *LightsAction[Message]) CallAction(storage storage.Storage, message Message, target *template.Button) {
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
