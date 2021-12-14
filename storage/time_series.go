package storage

import (
	"sync"
	"sync/atomic"
)

type bucketNode struct {
	ts    uint64
	next  *bucketNode
	value uint64
}

type timeSeries struct {
	layers []*tsLayer
}

type tsLayer struct {
	mu       sync.RWMutex
	first    *bucketNode
	nodes    *sync.Map //map[uint64]*bucketNode
	formatTs func(uint64) uint64
}

var (
	milliSecondsInMinute uint64 = 60000
	milliSecondsInHour   uint64 = 3600000
	milliSecondsInDay    uint64 = 86400000
	milliSecondsInWeek   uint64 = 604800000
	milliSecondsInMonth  uint64 = 2592000000
)

func newRootBucketNode() *bucketNode {
	return &bucketNode{
		ts:    0,
		value: 0,
		next:  nil,
	}
}

func newLayer(formatTs func(uint64) uint64) *tsLayer {
	return &tsLayer{
		first:    newRootBucketNode(),
		nodes:    &sync.Map{},
		formatTs: formatTs,
	}
}

func newTimeSeries() *timeSeries {
	layers := []*tsLayer{
		newLayer(tsToMonthBucket),
		newLayer(tsToWeekBucket),
		newLayer(tsToDayBucket),
		newLayer(tsToHourBucket),
		newLayer(tsToMinuteBucket),
	}
	return &timeSeries{
		layers: layers,
	}
}

func tsToMinuteBucket(ts uint64) uint64 {
	return ts/milliSecondsInDay + (ts%milliSecondsInDay)/milliSecondsInMinute
}

func tsToHourBucket(ts uint64) uint64 {
	return ts/milliSecondsInDay + (ts%milliSecondsInDay)/milliSecondsInHour
}

func tsToDayBucket(ts uint64) uint64 {
	return ts / milliSecondsInDay
}

func tsToWeekBucket(ts uint64) uint64 {
	return ts / milliSecondsInWeek
}

func tsToMonthBucket(ts uint64) uint64 {
	return ts / milliSecondsInMonth
}

func (t *timeSeries) add(ts uint64, value uint64) {
	for _, layer := range t.layers {
		layer.add(ts, value)
	}
}

func (l *tsLayer) add(ts uint64, value uint64) {
	tsFormatted := l.formatTs(ts)
	cachedNode, found := l.nodes.Load(tsFormatted)
	if !found {
		l.mu.Lock()
		cachedNode, found = l.nodes.Load(tsFormatted)
		if !found {
			node := l.findPrevBucketNode(tsFormatted)

			prevNext := node.next
			node.next = &bucketNode{
				ts:    ts,
				value: 0,
				next:  prevNext,
			}
			cachedNode = node.next
			l.nodes.Store(tsFormatted, cachedNode)
		}
		l.mu.Unlock()
	}
	atomic.AddUint64(&cachedNode.(*bucketNode).value, value)
}

func (t *timeSeries) getCount(startTs uint64, endTs uint64) uint64 {
	layerIdx := 0
	intervals := [][]uint64{{startTs, endTs}}
	var result uint64 = 0
	for len(intervals) > 0 && layerIdx < len(t.layers) {
		layer := t.layers[layerIdx]
		newIntervals := make([][]uint64, 0)
		for _, interval := range intervals {
			leftTs := layer.formatTs(interval[0])
			rightTs := layer.formatTs(interval[1])

			if rightTs > leftTs {
				result += layer.getCount(leftTs, rightTs)
				if leftTs-interval[0] > 0 {
					newIntervals = append(newIntervals, []uint64{interval[0], leftTs})
				}
			} else {
				newIntervals = append(newIntervals, []uint64{interval[0], interval[1]})
			}
		}
	}
	return result
}

func (l *tsLayer) getCount(startTs uint64, endTs uint64) uint64 {
	l.mu.RLock()
	defer l.mu.RUnlock()

	var sum uint64 = 0
	node := l.findPrevBucketNode(startTs).next
	for node != nil && node.ts <= endTs {
		sum += node.value
		node = node.next
	}
	return sum
}

func (l *tsLayer) findPrevBucketNode(ts uint64) *bucketNode {
	node := l.first
	for node.next != nil && node.next.ts < ts {
		node = node.next
	}
	return node
}
