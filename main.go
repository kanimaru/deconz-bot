package main

import (
	"github.com/PerformLine/go-stockutil/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	"telegram-deconz/view"
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
		deconzClient  = deconz.CreateDeconzClient(setting)
		deviceService = deconz.CreateDeconzDeviceService(deconzClient)

		apiKey                     = getEnv("TELEGRAM_API_KEY", "")
		api                        = telegram.CreateBot(apiKey)
		storageManager             = storage.CreateInMemoryStorage()
		actionManager              = telegram.CreateActionManager(storageManager, api.BotAPI)
		commandManager             = bot.CreateCommandManager[telegram.Message]()
		distributor                = bot.CreateMessageDistributor[telegram.Message]()
		viewOnClickHandler         = bot.CreateViewOnClickHandler[telegram.Message]()
		groupsOnClickHandler       = bot.CreateGroupsOnClickHandler[telegram.Message](deviceService, engine)
		lightsOnClickHandler       = bot.CreateLightsOnClickHandler[telegram.Message](deviceService, engine)
		lightsActionOnClickHandler = bot.CreateLightActionOnClickHandler[telegram.Message](deviceService, engine)
	)

	distributor.AddMessageReceiver(actionManager)
	distributor.AddMessageReceiver(commandManager)
	distributor.AddMessageReceiver(lightsActionOnClickHandler)
	actionManager.RegisterAction(viewOnClickHandler, "Action.Close", "Action.Back")
	actionManager.RegisterAction(groupsOnClickHandler, "Select.Group")
	actionManager.RegisterAction(lightsOnClickHandler, "Select.Light")
	actionManager.RegisterAction(lightsActionOnClickHandler, lightsActionOnClickHandler.HandledActions...)

	commandManager.AddCommand("light", bot.CommandDefinition[telegram.Message]{
		Description: "Let control lights and switches",
		Exec: func(message telegram.Message) {
			_, err := api.BotAPI.Request(tgbotapi.NewDeleteMessage(message.GetChatId(), message.GetId()))
			if err != nil {
				log.Warningf("Can't delete the request")
			}

			groupsMap := deviceService.GetGroups()
			groupsData := view.GroupsData{
				Groups: groupsMap,
			}

			groupView, err := engine.Apply("groups.go.xml", groupsData)
			if err != nil {
				panic(err)
			}

			viewManager := telegram.CreateViewManager(api.BotAPI)

			msg, err := viewManager.Open(groupView, message)
			log.Infof("Msg ID: %v", msg.GetId())
			s := storageManager.Get(msg.GetId())
			s.Save("viewManager", viewManager)
			if err != nil {
				log.Fatalf("Problem with changing keyboard: %w", err)
			}
		},
	})
	api.UpdateCommands(commandManager)

	go api.HandleUpdates(distributor.ReceiveMessage)
	<-doneChan
}
