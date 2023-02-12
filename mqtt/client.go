package mqtt

import (
	"github.com/PerformLine/go-stockutil/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func CreateMqttClient(options *mqtt.ClientOptions) mqtt.Client {
	mqttClient := mqtt.NewClient(options)
	token := mqttClient.Connect()
	token.WaitTimeout(5 * time.Second)
	err := token.Error()
	if err != nil {
		log.Fatalf("Can't connect MQTT: %w", err)
	}
	log.Infof("MQTT Connected")
	return mqttClient
}
