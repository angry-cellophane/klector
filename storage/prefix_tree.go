package storage

import (
	"sort"
	"sync"
)

type node struct {
	mu     sync.RWMutex
	series *sync.Map //map[string]*timeSeries
	nodes  *sync.Map
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
	child, found := n.nodes.Load(name)
	if !found {
		n.mu.Lock()
		child, found = n.nodes.Load(name)
		if !found {
			child = &node{
				series: &sync.Map{},
				nodes:  &sync.Map{},
			}
			n.nodes.Store(name, child)
		}
		n.mu.Unlock()
	}
	child.(*node).addToSeries(event.Timestamp, event.Attributes[name], 1)
	for i := 1; i < len(names); i++ {
		child.(*node).addChildNode(event, names[i:])
	}
}

func (n *node) addToSeries(ts uint64, attrValue string, count uint64) {
	series, found := n.series.Load(attrValue)
	if !found {
		n.mu.Lock()
		series, found = n.series.Load(attrValue)
		if !found {
			series = newTimeSeries()
			n.series.Store(attrValue, series)
		}
		n.mu.Unlock()
	}
	series.(*timeSeries).add(ts, count)
}

func findTimeSeries(n *node, names []string, query *Query) *timeSeries {
	if len(names) == 0 {
		return nil
	}

	child, found := n.nodes.Load(names[0])
	if !found {
		return nil
	}

	series, found := child.(*node).series.Load(query.Attributes[names[0]])
	if !found {
		return nil
	}
	if len(names) == 1 {
		return series.(*timeSeries)
	}

	return findTimeSeries(child.(*node), names[1:], query)
}

func (t *tree) find(query *Query) *timeSeries {
	names := sortAttributes(query.Attributes)
	return findTimeSeries(t.root, names, query)
}

func newTree() *tree {
	return &tree{
		root: &node{
			nodes: &sync.Map{},
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
