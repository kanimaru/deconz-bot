package template

type StorageManager interface {
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
}

type InMemoryStorageManager struct {
	// Map of messageId and Storage
	storages map[int]Storage
}

func CreateInMemoryStorage() *InMemoryStorageManager {
	return &InMemoryStorageManager{
		storages: make(map[int]Storage),
	}
}

func (i *InMemoryStorageManager) Add(messageId int) Storage {
	storage := &InMemoryStorage{
		data: make(map[string]interface{}),
	}
	i.storages[messageId] = storage
	return storage
}

func (i *InMemoryStorageManager) Get(messageId int) Storage {
	storage, ok := i.storages[messageId]
	if !ok {
		return i.Add(messageId)
	}
	return storage
}

func (i *InMemoryStorageManager) Remove(messageId int) {
	delete(i.storages, messageId)
}

type InMemoryStorage struct {
	// Map of date keys and their values
	data map[string]interface{}
}

func (i *InMemoryStorage) Save(key string, value interface{}) {
	i.data[key] = value
}

func (i *InMemoryStorage) Get(key string) interface{} {
	return i.data[key]
}

func (i *InMemoryStorage) Delete(key string) interface{} {
	value := i.data[key]
	delete(i.data, key)
	return value
}
