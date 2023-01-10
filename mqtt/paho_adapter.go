package mqtt

import (
	"github.com/PerformLine/go-stockutil/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type PahoClient struct {
	mqtt.Client
}

func CreateMqttClient(clientId string, username string, password string, urls ...string) Client {
	options := mqtt.NewClientOptions()
	options.Username = username
	options.Password = password
	options.ClientID = clientId
	options.AutoReconnect = true
	options.Store = mqtt.NewMemoryStore()
	for _, url := range urls {
		options.AddBroker(url)
	}
	mqttClient := mqtt.NewClient(options)
	mqttClient.Connect()
	log.Infof("MQTT Connected")
	return &PahoClient{
		Client: mqttClient,
	}
}

func (p *PahoClient) Send(topic string, message interface{}) {
	p.Client.Publish(topic, 0, false, message)
}
