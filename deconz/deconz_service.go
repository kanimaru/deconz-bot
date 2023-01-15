package deconz

import (
	"github.com/PerformLine/go-stockutil/log"
	"github.com/kanimaru/godeconz/http"
	"github.com/lucasb-eyer/go-colorful"
)

var on = true
var off = false

type LightState struct {
	// On optional will be set with brightness can be usefully when you want to hold the brightness
	On *bool
	// Brightness between 0-255
	Brightness *uint8
	// Color is an RGB hex string
	Color string
	// Temperature of the light between ctMin and ctMax of the light
	Temperature *int
}

type DeviceService interface {
	GetGroups() map[string]string
	GetLights() map[string]string
	GetLightsForGroup(group string) map[string]string
	GetLight(light string) http.LightResponseState
	GetGroup(group string) http.GroupResponseAttribute
	SetLightState(state LightState, lights ...string)
}

type deviceService[T any] struct {
	client *http.Client[T]
}

func CreateDeconzDeviceService[T any](client *http.Client[T]) DeviceService {
	log.Notice("Deconz device service initialized.")
	return deviceService[T]{
		client: client,
	}
}

func (d deviceService[T]) GetGroups() map[string]string {
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

func (d deviceService[T]) GetGroup(group string) http.GroupResponseAttribute {
	var groupResponse http.GroupResponseAttribute
	_, err := d.client.GetGroupAttributes(group, &groupResponse)
	if err != nil {
		log.Fatalf("Can't resolve group: %v", err)
	}
	return groupResponse

}

func (d deviceService[T]) GetLights() map[string]string {
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

func (d deviceService[T]) GetLight(light string) http.LightResponseState {
	var state http.LightResponseState
	_, err := d.client.Get("/lights/%s", &state, light)
	if err != nil {
		log.Fatalf("Can't resolve light: %v", err)
	}
	return state
}

func (d deviceService[T]) SetLightState(state LightState, lights ...string) {
	lightState := http.LightRequestState{
		Ct: state.Temperature,
		On: state.On,
	}

	if state.On == nil && state.Brightness != nil {
		if *state.Brightness > 0 {
			lightState.On = &on
		} else {
			lightState.On = &off
		}
	}

	if state.Color != "" {
		color, err := colorful.Hex(state.Color)
		if err != nil {
			log.Errorf("Can't convert %v to color: %v", state.Color, err)
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
		lightState.Bri = state.Brightness
	}

	for _, light := range lights {
		_, err := d.client.SetLightState(light, lightState)
		if err != nil {
			log.Fatalf("Can't update light %v: %v", light, err)
		}
	}
}

func (d deviceService[T]) GetLightsForGroup(group string) map[string]string {
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
