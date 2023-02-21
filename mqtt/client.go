package mqtt

import (
	"github.com/PerformLine/go-stockutil/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func CreateMqttClient(options *mqtt.ClientOptions) mqtt.Client {
	log.Infof("MQTT connecting...")
	mqttClient := mqtt.NewClient(options)
	token := mqttClient.Connect()
	if token.Wait() && token.Error() != nil {
		log.Fatalf("Can't connect MQTT: %w", token.Error())
	}
	log.Infof("MQTT connected!")
	return mqttClient
}
