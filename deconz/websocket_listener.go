package deconz

import (
	"fmt"
	"github.com/PerformLine/go-stockutil/log"
	"github.com/kanimaru/godeconz/ws"
)

type MessageSender interface {
	SendMessage(message string)
}

func ListenForAddedDevices(service Service, messageSender MessageSender) {
	for device := range service.GetAddedDevices() {
		var data ws.Attr
		if device.Light != nil {
			err := device.LightAs(&data)
			if err != nil {
				warnNotParseable()
			}
		} else if device.Sensor != nil {
			err := device.SensorAs(&data)
			if err != nil {
				warnNotParseable()
			}
		} else if device.Group != nil {
			// Could be buggy. Tested it but got no data from websocket so probably not working at all.
			err := device.GroupAs(&data)
			if err != nil {
				warnNotParseable()
			}
		}
		message := fmt.Sprintf("New %v was added called %v", data.Type, data.Name)
		messageSender.SendMessage(message)
	}
}

func warnNotParseable() {
	log.Warningf("Couldn't transform the data to Attr. To be honest it isn't the right data type " +
		"anyway but it's not documented.")
}
