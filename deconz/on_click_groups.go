package deconz

import (
	"github.com/kanimaru/godeconz/http"
	"telegram-deconz/bot"
	"telegram-deconz/storage"
	"telegram-deconz/template"
	"telegram-deconz/view"
)

type GroupsOnClickHandler[Message bot.BaseMessage] struct {
	deconzService Service
}

func CreateGroupsOnClickHandler[Message bot.BaseMessage](deconzService Service) *GroupsOnClickHandler[Message] {
	return &GroupsOnClickHandler[Message]{
		deconzService: deconzService,
	}
}

func (g *GroupsOnClickHandler[Message]) CallAction(storage storage.Storage, message Message, target *template.Button) {
	storage.Save("group", target.Data)
	views := storage.Get("viewManager").(bot.ViewManager[Message])

	lightMap := g.deconzService.GetLightsForGroup(target.Data)
	lights := make([]http.LightResponseState, 0, len(lightMap))
	for lightId := range lightMap {
		light := g.deconzService.GetLight(lightId)
		lights = append(lights, light)
	}

	features := GetLightFeatures(lights...)

	lightData := view.LightsData{
		GroupName:            target.Label,
		GroupId:              target.Data,
		Lights:               lightMap,
		On:                   features.On,
		ColorAvailable:       features.HasColor,
		BrightnessAvailable:  features.HasBrightness,
		TemperatureAvailable: features.HasTemp,
	}

	views.Show("lights.go.xml", lightData, message)
}
