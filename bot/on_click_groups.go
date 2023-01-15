package bot

import (
	"github.com/PerformLine/go-stockutil/log"
	"github.com/kanimaru/godeconz/http"
	"telegram-deconz/deconz"
	"telegram-deconz/storage"
	"telegram-deconz/template"
	"telegram-deconz/view"
)

type GroupsOnClickHandler[Message BaseMessage] struct {
	deconzService deconz.DeviceService
	engine        template.Engine
}

func CreateGroupsOnClickHandler[Message BaseMessage](deconzService deconz.DeviceService, engine template.Engine) *GroupsOnClickHandler[Message] {
	return &GroupsOnClickHandler[Message]{
		deconzService: deconzService,
		engine:        engine,
	}
}

func (g *GroupsOnClickHandler[Message]) CallAction(storage storage.Storage, message Message, target *template.Button) {
	storage.Save("group", target.Data)
	views := storage.Get("viewManager").(ViewManager[Message])

	lightMap := g.deconzService.GetLightsForGroup(target.Data)
	lights := make([]http.LightResponseState, 0, len(lightMap))
	for lightId := range lightMap {
		light := g.deconzService.GetLight(lightId)
		lights = append(lights, light)
	}

	features := deconz.GetLightFeatures(lights...)

	lightData := view.LightsData{
		GroupName:            target.Label,
		GroupId:              target.Data,
		Lights:               lightMap,
		On:                   features.On,
		ColorAvailable:       features.HasColor,
		BrightnessAvailable:  features.HasBrightness,
		TemperatureAvailable: features.HasTemp,
	}

	lightView, err := g.engine.Apply("lights.go.xml", lightData)
	if err != nil {
		_ = views.Close(message)
		log.Errorf("Problems with parsing template for lights - %v: %w", target.Label, err)
	}
	_, err = views.Open(lightView, message)
	if err != nil {
		_ = views.Close(message)
		log.Errorf("Problems with displaying lights for %v: %w", target.Label, err)
	}
}
