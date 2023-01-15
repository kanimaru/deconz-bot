package deconz

import (
	"github.com/PerformLine/go-stockutil/log"
	"telegram-deconz/bot"
	"telegram-deconz/storage"
	"telegram-deconz/template"
)

type ScanAction[Message bot.BaseMessage] struct {
	deconzService Service
}

func CreateScanAction[Message bot.BaseMessage](deconzService Service) *ScanAction[Message] {
	return &ScanAction[Message]{
		deconzService: deconzService,
	}
}

func (l *ScanAction[Message]) CallAction(storage storage.Storage, message Message, target *template.Button) {
	views := storage.Get("viewManager").(bot.ViewManager[Message])
	switch *target.OnClick {
	case "Action.StartScan":
		l.deconzService.StartScan(255)
		startResult := template.View{Text: "Scan started"}
		_, err := views.Open(&startResult, message)
		if err != nil {
			log.Errorf("Can't return scan success message")
		}
	case "Action.StopScan":
		l.deconzService.StopScan()
		stopResult := template.View{Text: "Scan stopped"}
		_, err := views.Open(&stopResult, message)
		if err != nil {
			log.Errorf("Can't return scan stop message")
		}
	}

}
