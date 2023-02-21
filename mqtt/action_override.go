package mqtt

import (
	"github.com/PerformLine/go-stockutil/log"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"telegram-deconz/bot"
	"telegram-deconz/storage"
	"telegram-deconz/template"
	"time"
)

type OverrideAction[Message bot.BaseMessage] struct {
	client        mqtt.Client
	readCallbacks map[string]func(overrideActive bool)
}

func CreateOverrideAction[Message bot.BaseMessage](client mqtt.Client) *OverrideAction[Message] {
	action := &OverrideAction[Message]{
		client:        client,
		readCallbacks: make(map[string]func(overrideActive bool)),
	}
	return action
}

func (g *OverrideAction[Message]) Call(storage storage.Storage, message Message, target *template.Button) {
	views := storage.Get("viewManager").(bot.ViewManager[Message])
	g.readCallbacks[target.Data] = func(overrideActive bool) {
		if overrideActive {
			view := template.View{Text: "Override activated for " + target.Data}
			_, _ = views.Open(&view, message)
		} else {
			view := template.View{Text: "Override deactivated for " + target.Data}
			_, _ = views.Open(&view, message)
		}
		delete(g.readCallbacks, target.Data)
	}

	mqttMessage := BaseMessage{
		From:    strconv.FormatInt(message.GetChatId(), 10),
		Payload: true,
	}
	publish := g.client.Publish(target.Data+"/override/write", 1, false, mqttMessage.ToJson())
	go func() {
		publish.WaitTimeout(5 * time.Second)
		err := publish.Error()

		var view *template.View
		if err == nil {
			view = &template.View{Text: "Send override request for " + target.Data}
		} else {
			view = &template.View{Text: "Override request failed for " + target.Data}
			log.Errorf("Can't send message to MQTT: %w", err)
		}
		_, err = views.Open(view, message)
		if err != nil {
			log.Errorf("Can't open success message for send override: %w", err)
		}
	}()
}
