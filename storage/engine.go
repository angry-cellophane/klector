package storage

type Event struct {
	Id         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
	Timestamp  uint64            `json:"timestamp"`
}

type Events struct {
	Events []Event `json:"events"`
}

type Query struct {
	Id             string            `json:"id"`
	Attributes     map[string]string `json:"attributes"`
	StartTimestamp uint64            `json:"startTimestamp"`
	EndTimestamp   uint64            `json:"endTimestamp"`
}

type ResultSet struct {
	Id         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
	Value      uint64            `json:"value"`
}

type Storage interface {
	Write(events *Events) error
	Query(query *Query) (*ResultSet, error)
}

type StorageConfiguration struct {
	DataFolder string `json:"dataFolder"`
}

func NewDefaultStorageConfiguration() *StorageConfiguration {
	return &StorageConfiguration{
		DataFolder: "./data",
	}
}

func Create(config *StorageConfiguration) Storage {
	return &inMemoryStorage{
		tree: newTree(),
	}
}
