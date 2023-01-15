package bot

import "telegram-deconz/template"

// ViewManager is a helper to switch between different views.
type ViewManager[Message BaseMessage] interface {
	// Open a new view and pushes the old one on the stack
	Open(view *template.View, message Message) (Message, error)
	// Back to the previous view
	Back(message Message) (*template.View, error)
	// Close the current message / inline context
	Close(message Message) error
	// FindButton locates a button in the current view
	FindButton(id string) *template.Button
	// Show a template with given data
	Show(name string, data interface{}, message Message) (*template.View, *Message)
}
