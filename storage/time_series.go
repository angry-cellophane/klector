package storage

import "sync"

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

func newTimeSeries() *timeSeries {
	return &timeSeries{
		first: &bucketNode{
			ts:    0,
			value: 0,
			next:  nil,
		},
		nodes: make(map[uint64]*bucketNode),
	}
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
