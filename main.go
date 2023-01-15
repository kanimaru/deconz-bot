package main

import (
	"github.com/PerformLine/go-stockutil/log"
	"github.com/kanimaru/godeconz"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"telegram-deconz/bot"
	"telegram-deconz/deconz"
	"telegram-deconz/storage"
	"telegram-deconz/telegram"
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
	log.Infof("ending app")
	doneChan <- true
}

func main() {
	// Rand seed is needed to generate the same names for multiple runs
	rand.Seed(1337)
	doneChan = make(chan bool)
	go handleSignals()

	var (
		engine  = template.CreateEngineByDir("view/")
		setting = godeconz.Settings{
			Address:      getEnv("DECONZ_ADDRESS", ""),
			HttpProtocol: getEnv("DECONZ_PROTO", "http"),
			ApiKey:       getEnv("DECONZ_API_KEY", ""),
		}
		deconzClient  = deconz.CreateClient(setting)
		deconzService = deconz.CreateService(deconzClient)

		apiKey                     = getEnv("TELEGRAM_API_KEY", "")
		tgBot                      = telegram.CreateBot(apiKey)
		storageManager             = storage.CreateInMemoryStorage()
		commands                   = telegram.CreateCommand(tgBot, deconzService, storageManager, engine)
		actionManager              = telegram.CreateActionManager(storageManager)
		commandManager             = bot.CreateCommandManager[telegram.Message]()
		distributor                = bot.CreateMessageDistributor[telegram.Message]()
		viewOnClickHandler         = bot.CreateViewOnClickHandler[telegram.Message]()
		groupsOnClickHandler       = deconz.CreateGroupsOnClickHandler[telegram.Message](deconzService)
		lightsOnClickHandler       = deconz.CreateLightsOnClickHandler[telegram.Message](deconzService)
		lightsActionOnClickHandler = deconz.CreateLightActionOnClickHandler[telegram.Message](deconzService)
	)

	distributor.AddMessageReceiver(actionManager)
	distributor.AddMessageReceiver(commandManager)
	distributor.AddMessageReceiver(lightsActionOnClickHandler)
	actionManager.RegisterAction(viewOnClickHandler, "Action.Close", "Action.Back")
	actionManager.RegisterAction(groupsOnClickHandler, "Select.Group")
	actionManager.RegisterAction(lightsOnClickHandler, "Select.Light")
	actionManager.RegisterAction(lightsActionOnClickHandler, lightsActionOnClickHandler.HandledActions...)
	commandManager.AddCommand("light", commands.CreateLightCmd())

	tgBot.UpdateCommands(commandManager)

	go tgBot.HandleUpdates(distributor.ReceiveMessage)
	<-doneChan
}
