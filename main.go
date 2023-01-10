package main

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"os/signal"
	"syscall"
	"telegram-deconz/template"
)

func getEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	if len(fallback) == 0 {
		log.Fatalf("Missing %q", key)
	}
	return fallback
}

var doneChan chan bool

func handleSignals() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-c
	doneChan <- true
}

func main() {
	go handleSignals()

	groups := GroupsData{
		Groups: []Group{
			{
				Name: "Group1",
				Data: "Data1",
			},
			{
				Name: "Group2",
				Data: "Data2",
			},
		},
	}
	engine := template.CreateEngineByDir("view/")
	apply, err := engine.Apply("groups.yaml", groups)
	if err != nil {
		panic(err)
	}
	log.Infof("%+v", apply)

	deconzClient := createDeconzClient()
	deviceService := createDeconzDeviceService(deconzClient)
	telegramClient := createTelegramClient()

	data := Data{}
	lightSetting := asView(createLightSettingCallback(&data, deviceService, telegramClient), telegramClient)
	listLight := asView(
		allLights(
			createLightContext(&data, "Select lights you want to change:", deviceService, lightSetting),
			&data,
			deviceService,
			telegramClient,
		),
		telegramClient)

	listGroup := asView(createGroupContext(&data, "Select the group you want to change:", deviceService, listLight), telegramClient)

	telegramClient.AddCommand("light", CommandCallback{
		description: "Controls the light",
		exec: func(update tgbotapi.Update) {
			telegramClient.changeContext(update, listGroup)
		},
	})
	telegramClient.UpdateCommands()
	telegramClient.handleUpdates(doneChan)
}

type Data struct {
	selectedGroup string
	selectedLight string
}

func (d *Data) GetLights() []string {
	return []string{d.selectedLight}
}

func (d *Data) GetLight() string {
	return d.selectedLight
}

func (d *Data) GetGroup() string {
	return d.selectedGroup
}

func (d *Data) SetLight(lightName string) {
	d.selectedLight = lightName
}

func (d *Data) SetGroup(group string) {
	d.selectedGroup = group
}
