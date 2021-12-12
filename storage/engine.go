package storage

type Event struct {
	Id         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
	Timestamp  uint64            `json:"timestamp"`
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
	Write(event *Event) error
	Query(query *Query) *ResultSet
}

func Create() Storage {
	return &inMemoryStorage{
		attributes: make(map[string]map[string]*timeSeries),
	}
}
