package storage

type Event struct {
	Id         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

type Query struct {
	Id         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
}

type ResultSet struct {
	Id         string            `json:"id"`
	Attributes map[string]string `json:"attributes"`
	Value      int64             `json:"value"`
}

type Storage interface {
	Write(event *Event)
	Query(query *Query) *ResultSet
}

type inMemoryStorage struct {
}

func (s *inMemoryStorage) Write(event *Event) {
	// noop
}
func (s *inMemoryStorage) Query(query *Query) *ResultSet {
	// noop
	return &ResultSet{}
}

func Create() Storage {
	return &inMemoryStorage{}
}
