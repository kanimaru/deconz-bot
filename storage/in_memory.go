package storage

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

func (i *InMemoryStorage) Has(key string) bool {
	_, ok := i.data[key]
	return ok
}
