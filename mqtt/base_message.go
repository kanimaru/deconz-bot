package mqtt

import (
	"encoding/json"
	"github.com/PerformLine/go-stockutil/log"
)

type BaseMessage struct {
	From    string      `json:"from"`
	To      string      `json:"to,omitempty"`
	Payload interface{} `json:"payload"`
}

func (b BaseMessage) ToJson() []byte {
	jsonValue, err := json.Marshal(b)
	if err != nil {
		log.Errorf("Problem with marshaling BaseMessage: %w", err)
		return []byte{}
	}
	return jsonValue
}
