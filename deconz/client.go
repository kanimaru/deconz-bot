package deconz

import (
	"github.com/PerformLine/go-stockutil/log"
	"github.com/go-resty/resty/v2"
	"github.com/gorilla/websocket"
	"github.com/kanimaru/godeconz"
	"github.com/kanimaru/godeconz/http"
	"github.com/kanimaru/godeconz/ws"
)

func CreateHttpClient(setting godeconz.Settings) *http.Client[*resty.Response] {
	httpAdapter := http.CreateAdapterHttpClientResty(resty.New(), Logger{}, false)
	deconzClient := http.CreateClient(httpAdapter, setting)
	log.Notice("Deconz initialized.")
	return &deconzClient
}

func CreateWsClient(httpClient *http.Client[*resty.Response]) *ws.Client {
	logger := Logger{}
	deconzWsAdapter := ws.CreateAdapterWebsocketClientGorilla(websocket.DefaultDialer, logger)
	return ws.CreateClientFromConfig[*resty.Response](*httpClient, deconzWsAdapter, logger)
}

type LightFeatures struct {
	On            bool
	HasBrightness bool
	HasColor      bool
	HasTemp       bool
	Reachable     bool
}

func GetLightFeatures(lights ...http.LightResponseState) LightFeatures {
	features := LightFeatures{}
	for _, light := range lights {
		features.Reachable = features.Reachable || (light.State.Reachable != nil && *light.State.Reachable)
		features.HasTemp = features.HasTemp || (light.Ctmin != nil && light.Ctmax != nil && (*light.Ctmin != *light.Ctmax))
		features.HasColor = features.HasColor || (light.State.Colormode != nil && (*light.State.Colormode == http.ColorModeHS || *light.State.Colormode == http.ColorModeXY))
		features.On = features.On || (light.State.On != nil && *light.State.On)
		features.HasBrightness = features.HasBrightness || light.State.Bri != nil
	}
	return features
}
