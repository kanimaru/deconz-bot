package mqtt

type Client interface {
	// Send message to topic will be marshalled depending on adapter.
	Send(topic string, message interface{})
}
