package main

import (
	"github.com/PerformLine/go-stockutil/log"
	"github.com/go-resty/resty/v2"
	"github.com/kanimaru/godeconz"
	"github.com/kanimaru/godeconz/http"
)

func createDeconzClient() *http.Client[*resty.Response] {
	setting := godeconz.Settings{
		Address:      getEnv("DECONZ_ADDRESS", ""),
		HttpProtocol: getEnv("DECONZ_PROTO", "http"),
		ApiKey:       getEnv("DECONZ_API_KEY", ""),
	}

	httpAdapter := http.CreateAdapterHttpClientResty(resty.New(), Logger{}, false)
	deconzClient := http.CreateClient(httpAdapter, setting)
	log.Notice("Deconz initialized.")
	return &deconzClient
}
