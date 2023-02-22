package main

import (
	"encoding/json"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"telegram-deconz/bot"
	"telegram-deconz/deconz"
	"telegram-deconz/mqtt"
	"telegram-deconz/storage"
	"telegram-deconz/telegram"
	"telegram-deconz/template"
	"time"

	"github.com/PerformLine/go-stockutil/log"
	mqtt2 "github.com/eclipse/paho.mqtt.golang"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kanimaru/godeconz"
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
		engine             = template.CreateEngineByDir("view/")
		deconzHttpClient   = deconz.CreateHttpClient(getDeconzOptions())
		deconzWsClient     = deconz.CreateWsClient(deconzHttpClient)
		deconzService      = deconz.CreateService(deconzHttpClient, deconzWsClient)
		mqttClient         = mqtt.CreateMqttClient(getMqttOptions())
		chatId             = getChatId()
		apiKey             = getEnv("TELEGRAM_API_KEY", "")
		tgBot              = telegram.CreateBot(apiKey)
		basicMessageSender = CreateBaseMessageSender(chatId, tgBot)
		storageManager     = storage.CreateInMemoryStorage()
		commandFactory     = telegram.CreateCommandFactory(tgBot, deconzService, storageManager, engine)
		actionManager      = telegram.CreateActionManager(storageManager)
		commandManager     = bot.CreateCommandManager[telegram.Message]()
		distributor        = bot.CreateMessageDistributor[telegram.Message](storageManager)
		viewAction         = bot.CreateViewAction[telegram.Message]()
		groupsAction       = deconz.CreateGroupsAction[telegram.Message](deconzService)
		lightsAction       = deconz.CreateLightsAction[telegram.Message](deconzService)
		lightAction        = deconz.CreateLightAction[telegram.Message](deconzService)
		scanAction         = deconz.CreateScanAction[telegram.Message](deconzService)
		overrideAction     = mqtt.CreateOverrideAction[telegram.Message](mqttClient)
	)

	distributor.AddMessageReceiver(actionManager)
	distributor.AddMessageReceiver(commandManager)
	distributor.AddMessageReceiver(lightAction)
	actionManager.RegisterAction(viewAction, "Action.Close", "Action.Back")
	actionManager.RegisterAction(groupsAction, "Select.Group")
	actionManager.RegisterAction(lightsAction, "Select.Light")
	actionManager.RegisterAction(lightAction, lightAction.HandledActions...)
	actionManager.RegisterAction(scanAction, "Action.StartScan", "Action.StopScan")
	actionManager.RegisterAction(overrideAction, "Action.Override")
	commandManager.AddCommand("light", commandFactory.CreateLightCmd())
	commandManager.AddCommand("scan", commandFactory.CreateScanCmd())
	tgBot.UpdateCommands(commandManager, tgbotapi.NewBotCommandScopeChat(chatId))

	go deconzWsClient.Connect()
	go deconz.ListenForAddedDevices(deconzService, basicMessageSender)
	listenForChat(mqttClient, tgBot)
	go tgBot.HandleUpdates(distributor.ReceiveMessage)
	log.Infof("Init is ready! Start working...")
	<-doneChan
}

func getDeconzOptions() godeconz.Settings {
	return godeconz.Settings{
		Address:      getEnv("DECONZ_ADDRESS", ""),
		HttpProtocol: getEnv("DECONZ_PROTO", "http"),
		ApiKey:       getEnv("DECONZ_API_KEY", ""),
	}
}

func getMqttOptions() *mqtt2.ClientOptions {
	options := mqtt2.NewClientOptions().
		SetUsername(getEnv("MQTT_USERNAME", "")).
		SetPassword(getEnv("MQTT_PASSWORD", "")).
		SetClientID(getEnv("MQTT_CLIENT_ID", "deconzBot")).
		SetAutoReconnect(true).
		SetStore(mqtt2.NewMemoryStore()).
		SetPingTimeout(10 * time.Second).
		SetKeepAlive(10 * time.Second).
		SetResumeSubs(true).
		SetCleanSession(true).
		SetConnectionLostHandler(func(client mqtt2.Client, err error) {
			log.Errorf("Connection to MQTT broker lost.")
		}).
		SetOnConnectHandler(func(client mqtt2.Client) {
			log.Infof("Connection to MQTT broker established")
		})

	urls := getEnv("MQTT_URL", "")
	urlSlice := strings.Split(urls, "|")
	for _, broker := range urlSlice {
		options.AddBroker(broker)
	}
	return options
}

func getChatId() int64 {
	chatIdStr := getEnv("TELEGRAM_CHAT_ID", "")
	chatId, err := strconv.ParseInt(chatIdStr, 10, 64)
	if err != nil {
		log.Fatalf("Can't get telegram scope.")
	}
	return chatId
}

func listenForChat(client mqtt2.Client, bot telegram.Bot) {
	token := client.Subscribe("global/chat", 1, func(client mqtt2.Client, message mqtt2.Message) {
		var baseMessage mqtt.BaseMessage
		err := json.Unmarshal(message.Payload(), &baseMessage)
		if err != nil {
			log.Errorf("Can't unmarshal message lost telegram message: %w", err)
			return
		}

		chatId, err := strconv.ParseInt(baseMessage.To, 10, 64)
		if err != nil {
			log.Errorf("Can't parse chatId from %v -> %v: %v", baseMessage.From, baseMessage.To, err)
			return
		}

		msg := tgbotapi.NewMessage(chatId, baseMessage.Payload.(string))
		_, err = bot.Send(msg)
		if err != nil {
			log.Errorf("Can't send message received from MQTT: %v", err)
			return
		}
	})
	if token.Wait() && token.Error() != nil {
		log.Errorf("Couldn't subscribe to chat messages: %v", token.Error())
	}
	log.Infof("Listen for 'global/chat'")
}
