package bot

import (
	"github.com/PerformLine/go-stockutil/log"
	"telegram-deconz/storage"
	"telegram-deconz/template"
)

type ViewAction[Message BaseMessage] struct {
}

func CreateViewAction[Message BaseMessage]() *ViewAction[Message] {
	return &ViewAction[Message]{}
}

func (v *ViewAction[Message]) Call(storage storage.Storage, message Message, target *template.Button) {
	views := storage.Get("viewManager").(ViewManager[Message])

	switch *target.OnClick {
	case "Action.Close":
		err := views.Close(message)
		if err != nil {
			log.Errorf("Can't close view: %v", err)
		}
	case "Action.Back":
		prevView, err := views.Back(message)
		if err != nil {
			log.Errorf("Can't use previous view %v: %v", prevView.Name, err)
		}
	}
}
