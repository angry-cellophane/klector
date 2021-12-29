package storage

import (
	"sort"
	"sync"
)

type node struct {
	mu                 sync.RWMutex
	tseriesByAttrValue *sync.Map //map[string]*timeSeries where string is attribute value
	childNodes         *sync.Map //map[string]*node where string is attribute key
}

type tree struct {
	root *node
}

func (t *tree) addEvent(event *Event) {
	names := sortAttributes(event.Attributes)

	for len(names) > 0 {
		t.root.addChildNode(event, names)
		names = names[1:]
	}
}

func (n *node) addChildNode(event *Event, names []string) {
	if len(names) == 0 {
		return
	}

	name := names[0]
	child, found := n.childNodes.Load(name)
	if !found {
		n.mu.Lock()
		child, found = n.childNodes.Load(name)
		if !found {
			child = &node{
				tseriesByAttrValue: &sync.Map{},
				childNodes:         &sync.Map{},
			}
			n.childNodes.Store(name, child)
		}
		n.mu.Unlock()
	}
	child.(*node).addToSeries(event.Timestamp, event.Attributes[name], 1)
	for i := 1; i < len(names); i++ {
		child.(*node).addChildNode(event, names[i:])
	}
}

func (n *node) addToSeries(ts uint64, attrValue string, count uint64) {
	series, found := n.tseriesByAttrValue.Load(attrValue)
	if !found {
		n.mu.Lock()
		series, found = n.tseriesByAttrValue.Load(attrValue)
		if !found {
			series = newTimeSeries()
			n.tseriesByAttrValue.Store(attrValue, series)
		}
		n.mu.Unlock()
	}
	series.(*timeSeriesAggregator).add(ts, count)
}

func findTimeSeries(n *node, names []string, query *Query) *timeSeriesAggregator {
	if len(names) == 0 {
		return nil
	}

	child, found := n.childNodes.Load(names[0])
	if !found {
		return nil
	}

	series, found := child.(*node).tseriesByAttrValue.Load(query.Attributes[names[0]])
	if !found {
		return nil
	}
	if len(names) == 1 {
		return series.(*timeSeriesAggregator)
	}

	return findTimeSeries(child.(*node), names[1:], query)
}

func (t *tree) find(query *Query) *timeSeriesAggregator {
	names := sortAttributes(query.Attributes)
	return findTimeSeries(t.root, names, query)
}

func newTree() *tree {
	return &tree{
		root: &node{
			childNodes: &sync.Map{},
		},
	}
}

func sortAttributes(attributes map[string]string) []string {
	i := 0
	names := make([]string, len(attributes))
	for name := range attributes {
		names[i] = name
		i++
	}
	sort.Strings(names)
	return names
}
