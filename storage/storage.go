package storage

type Manager interface {
	// Add a storage for the given message id
	Add(messageId int) Storage
	// Get the storage for the given message id
	Get(messageId int) Storage
	// Remove the storage for given message id
	Remove(messageId int)
}

type Storage interface {
	// Save given key with value
	Save(key string, value interface{})
	// Get the data of the given key
	Get(key string) interface{}
	// Delete the given key returns the old value or nil if not exists
	Delete(key string) interface{}
	// Has checks if the key is available
	Has(key string) bool
}
