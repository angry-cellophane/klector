package storage

import (
	"errors"
	"log"
	"sync"
)

type bucketNode struct {
	ts    uint64
	next  *bucketNode
	value uint64
}

type timeSeries struct {
	mu    sync.RWMutex
	first *bucketNode
	nodes map[uint64]*bucketNode
}

func (t *timeSeries) add(ts uint64, value uint64) {
	t.mu.Lock()
	defer t.mu.Unlock()

	node, found := t.nodes[ts]
	if !found {
		node = t.findPrevBucketNode(ts)

		prevNext := node.next
		node.next = &bucketNode{
			ts:    ts,
			value: 0,
			next:  prevNext,
		}
		node = node.next
		t.nodes[ts] = node
	}
	node.value = node.value + value
}

func (t *timeSeries) getCount(startTs uint64, endTs uint64) uint64 {
	t.mu.RLock()
	defer t.mu.RUnlock()

	var sum uint64 = 0
	node := t.findPrevBucketNode(startTs).next
	for node != nil && node.ts <= endTs {
		sum += node.value
		node = node.next
	}
	return sum
}

func (t *timeSeries) findPrevBucketNode(ts uint64) *bucketNode {
	node := t.first
	for node.next != nil && node.next.ts < ts {
		node = node.next
	}
	return node
}

type inMemoryStorage struct {
	mu         sync.RWMutex
	attributes *sync.Map // map[string]map[string]*timeSeries
}

func (s *inMemoryStorage) Write(event *Event) error {
	if len(event.Attributes) == 0 {
		return errors.New("attributes are not defined in event")
	}
	if event.Timestamp == 0 {
		return errors.New("timestamp cannot be 0")
	}

	for name, value := range event.Attributes {
		values, found := s.attributes.Load(name)
		if !found {
			s.mu.Lock()
			_values := &sync.Map{} //make(map[string]*timeSeries)
			_values.Store(value, &timeSeries{
				first: &bucketNode{
					ts:    0,
					value: 0,
					next:  nil,
				},
				nodes: make(map[uint64]*bucketNode),
			})
			values = _values
			s.attributes.Store(name, values)
			s.mu.Unlock()
		}
		series, found := values.(*sync.Map).Load(value)
		if !found {
			s.mu.Lock()
			series = &timeSeries{
				first: &bucketNode{
					ts:    0,
					value: 0,
					next:  nil,
				},
				nodes: make(map[uint64]*bucketNode),
			}
			values.(*sync.Map).Store(value, series)
			s.mu.Unlock()
		}
		series.(*timeSeries).add(event.Timestamp, 1)
	}
	log.Printf("Received event %v", *event)

	return nil
}

func (s *inMemoryStorage) Query(query *Query) *ResultSet {
	var count uint64 = 0
	for name, value := range query.Attributes {
		attr, found := s.attributes.Load(name)
		if !found {
			continue
		}

		series, found := attr.(*sync.Map).Load(value)
		if !found {
			continue
		}

		count += series.(*timeSeries).getCount(query.StartTimestamp, query.EndTimestamp)
	}

	return &ResultSet{
		Id:         query.Id,
		Attributes: query.Attributes,
		Value:      count,
	}
}
