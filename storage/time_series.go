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

type timeSeriesAggregator struct {
	name     string
	mu       sync.RWMutex
	first    *bucketNode
	nodes    *sync.Map           //map[uint64]*bucketNode
	formatTs func(uint64) uint64 // format ts to bucket ts, assumes bucket ts <= input ts
	timeStep uint64
	subRange *timeSeriesAggregator
}

func newRootBucketNode() *bucketNode {
	return &bucketNode{
		ts:    0,
		value: 0,
		next:  nil,
	}
}

func newTimeSeries() *timeSeriesAggregator {
	return newTimeSeriesWithFormatter("month", milliSecondsInMonth, tsToMonthBucket,
		newTimeSeriesWithFormatter("day", milliSecondsInDay, tsToDayBucket,
			newTimeSeriesWithFormatter("hour", milliSecondsInHour, tsToHourBucket,
				newTimeSeriesWithFormatter("minute", milliSecondsInMinute, tsToMinuteBucket,
					nil))))
}

func newTimeSeriesWithFormatter(name string,
	timeStep uint64,
	formatTs func(uint64) uint64,
	subRange *timeSeriesAggregator) *timeSeriesAggregator {
	return &timeSeriesAggregator{
		name:     name,
		first:    newRootBucketNode(),
		nodes:    &sync.Map{},
		timeStep: timeStep,
		formatTs: formatTs,
		subRange: subRange,
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

func tsToMonthBucket(ts uint64) uint64 {
	return ts / milliSecondsInMonth
}

func (aggregator *timeSeriesAggregator) add(ts uint64, value uint64) {
	tsFormatted := aggregator.formatTs(ts)
	cachedNode, found := aggregator.nodes.Load(tsFormatted)
	if !found {
		aggregator.mu.Lock()
		cachedNode, found = aggregator.nodes.Load(tsFormatted)
		if !found {
			node := aggregator.findPrevBucketNode(tsFormatted)

			prevNext := node.next
			node.next = &bucketNode{
				ts:    tsFormatted,
				value: 0,
				next:  prevNext,
			}
			cachedNode = node.next
			aggregator.nodes.Store(tsFormatted, cachedNode)
		}
		aggregator.mu.Unlock()
	}
	atomic.AddUint64(&cachedNode.(*bucketNode).value, value)
	if aggregator.subRange != nil {
		aggregator.subRange.add(ts, value)
	}
}

func (aggregator *timeSeriesAggregator) getCount(startTs uint64, endTs uint64) uint64 {
	var leftBucketTs uint64 = aggregator.formatTs(startTs)
	var rightBucketTs uint64 = aggregator.formatTs(endTs)

	if aggregator.subRange != nil && rightBucketTs-leftBucketTs < aggregator.timeStep {
		return aggregator.subRange.getCount(startTs, endTs)
	}

	var result uint64 = aggregator.countInBuckets(leftBucketTs, rightBucketTs)
	if aggregator.subRange != nil {
		if leftBucketTs < startTs {
			result -= aggregator.subRange.getCount(leftBucketTs, startTs)
		}

		if rightBucketTs < endTs {
			result += aggregator.subRange.getCount(rightBucketTs, endTs)
		}
	}
	return result
}

/**
	Calculate sum of all buckets in range, endBucket is exclusive.
**/
func (aggregator *timeSeriesAggregator) countInBuckets(startBucket uint64, endBucket uint64) uint64 {
	aggregator.mu.RLock()
	defer aggregator.mu.RUnlock()

	var sum uint64 = 0
	node := aggregator.findPrevBucketNode(startBucket).next
	for node != nil && node.ts <= endBucket {
		sum += node.value
		node = node.next
	}
	return sum
}

func (aggregator *timeSeriesAggregator) findPrevBucketNode(ts uint64) *bucketNode {
	node := aggregator.first
	for node.next != nil && node.next.ts < ts {
		node = node.next
	}
	return node
}
