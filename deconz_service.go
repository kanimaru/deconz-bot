package main

import (
	"github.com/PerformLine/go-stockutil/log"
	"github.com/kanimaru/godeconz/http"
	"github.com/lucasb-eyer/go-colorful"
)

type LightState struct {
	// on optional will be set with brightness can be usefully when you want to hold the brightness
	on *bool
	// brightness between 0-255
	brightness *uint8
	// color is an RGB hex string
	color string
	// temperature of the light between ctMin and ctMax of the light
	temperature *int
}

type DeviceService interface {
	GetGroups() map[string]string
	GetLights() map[string]string
	GetLightsForGroup(group string) map[string]string
	GetLight(light string) http.LightResponseState
	GetGroup(group string) http.GroupResponseAttribute
	SetLightState(state LightState, lights ...string)
}

type DeconzDeviceService[T any] struct {
	client *http.Client[T]
}

func createDeconzDeviceService[T any](client *http.Client[T]) DeviceService {
	log.Notice("Deconz device service initialized.")
	return DeconzDeviceService[T]{
		client: client,
	}
}

func (d DeconzDeviceService[T]) GetGroups() map[string]string {
	groups := make(map[string]http.GroupResponse)
	_, err := d.client.GetAllGroups(&groups)
	if err != nil {
		log.Fatalf("Can't resolve groups: %v", err)
	}
	groupNames := make(map[string]string)
	for id, group := range groups {
		groupNames[id] = group.Name
	}
	return groupNames
}

func (d DeconzDeviceService[T]) GetGroup(group string) http.GroupResponseAttribute {
	var groupResponse http.GroupResponseAttribute
	_, err := d.client.GetGroupAttributes(group, &groupResponse)
	if err != nil {
		log.Fatalf("Can't resolve group: %v", err)
	}
	return groupResponse

}

func (d DeconzDeviceService[T]) GetLights() map[string]string {
	lights := make(map[string]http.LightResponseState)
	_, err := d.client.GetAllLights(&lights)
	if err != nil {
		log.Fatalf("Can't resolve lights: %v", err)
	}
	lightNames := make(map[string]string)
	for id, light := range lights {
		lightNames[id] = light.Name
	}
	return lightNames
}

func (d DeconzDeviceService[T]) GetLight(light string) http.LightResponseState {
	var state http.LightResponseState
	_, err := d.client.Get("/lights/%s", &state, light)
	if err != nil {
		log.Fatalf("Can't resolve light: %v", err)
	}
	return state
}

func (d DeconzDeviceService[T]) SetLightState(state LightState, lights ...string) {
	lightState := http.LightRequestState{
		Ct: state.temperature,
		On: state.on,
	}

	if state.on == nil && state.brightness != nil {
		if *state.brightness > 0 {
			lightState.On = &on
		} else {
			lightState.On = &off
		}
	}

	if state.color != "" {
		color, err := colorful.Hex(state.color)
		if err != nil {
			log.Errorf("Can't convert %v to color: %v", state.color, err)
		}

		h, s, v := color.Hsv()

		v *= 255
		vpt := uint8(v)

		s *= 255
		spt := uint8(s)

		h = (float64(h) / 360) * 65535
		hpt := uint16(h)

		lightState.Bri = &vpt
		lightState.Hue = &hpt
		lightState.Sat = &spt
	} else {
		lightState.Bri = state.brightness
	}

	for _, light := range lights {
		_, err := d.client.SetLightState(light, lightState)
		if err != nil {
			log.Fatalf("Can't update light %v: %v", light, err)
		}
	}
}

func (d DeconzDeviceService[T]) GetLightsForGroup(group string) map[string]string {
	var groupResponse http.GroupResponseAttribute
	_, err := d.client.GetGroupAttributes(group, &groupResponse)
	if err != nil {
		log.Fatalf("Can't resolve group %v: %v", group, err)
	}
	lightNames := make(map[string]string)
	for _, lightId := range groupResponse.Lights {

		var lightState http.LightResponseState
		_, err = d.client.GetLightState(lightId, &lightState)
		if err != nil {
			log.Fatalf("Can't resolve lights: %v", err)
		}
		lightNames[lightId] = lightState.Name
	}
	return lightNames
}
