package storage

import (
	"errors"
	"log"
)

type inMemoryStorage struct {
	tree            *tree
	formatTimestamp func(timestamp uint64) uint64
}

func (s *inMemoryStorage) Write(events *Events) error {
	for _, event := range events.Events {
		if err := s.writeEvent(&event); err != nil {
			return err
		}
	}
	return nil
}

func (s *inMemoryStorage) writeEvent(event *Event) error {
	if len(event.Attributes) == 0 {
		return errors.New("attributes are not defined in event")
	}
	if event.Timestamp == 0 {
		return errors.New("timestamp cannot be 0")
	}

	log.Printf("Received event %v", *event)
	event.Timestamp = s.formatTimestamp(event.Timestamp)
	s.tree.addEvent(event)

	return nil
}

func (s *inMemoryStorage) Query(query *Query) (*ResultSet, error) {
	var count uint64 = 0

	series := s.tree.find(query)
	if series != nil {
		count = series.getCount(query.StartTimestamp, query.EndTimestamp)
	}

	return &ResultSet{
		Id:         query.Id,
		Attributes: query.Attributes,
		Value:      count,
	}, nil
}
