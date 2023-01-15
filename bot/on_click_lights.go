package bot

import (
	"github.com/PerformLine/go-stockutil/log"
	"telegram-deconz/deconz"
	"telegram-deconz/storage"
	"telegram-deconz/template"
	"telegram-deconz/view"
)

type LightsOnClickHandler[Message BaseMessage] struct {
	deconzService deconz.DeviceService
	engine        template.Engine
}

func CreateLightsOnClickHandler[Message BaseMessage](deconzService deconz.DeviceService, engine template.Engine) *LightsOnClickHandler[Message] {
	return &LightsOnClickHandler[Message]{
		deconzService: deconzService,
		engine:        engine,
	}
}

func (l *LightsOnClickHandler[Message]) CallAction(storage storage.Storage, message Message, target *template.Button) {
	storage.Save("light", target.Data)
	groupId := storage.Get("group").(string)
	group := l.deconzService.GetGroup(groupId)
	views := storage.Get("viewManager").(ViewManager[Message])

	light := l.deconzService.GetLight(target.Data)
	features := deconz.GetLightFeatures(light)

	lightData := view.LightData{
		GroupName:            group.Name,
		Id:                   target.Data,
		Name:                 light.Name,
		On:                   features.On,
		ColorAvailable:       features.HasColor,
		BrightnessAvailable:  features.HasBrightness,
		TemperatureAvailable: features.HasTemp,
	}

	lightView, err := l.engine.Apply("light.go.xml", lightData)
	if err != nil {
		_ = views.Close(message)
		log.Errorf("Problems with parsing template for light - %v: %w", target.Label, err)
	}
	_, err = views.Open(lightView, message)
	if err != nil {
		_ = views.Close(message)
		log.Errorf("Problems with displaying light for %v: %w", target.Label, err)
	}
}
