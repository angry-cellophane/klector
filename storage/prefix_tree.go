package storage

import (
	"sort"
	"sync"
)

type node struct {
	mu     sync.RWMutex
	series *timeSeries
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
				series: newTimeSeries(),
				nodes:  &sync.Map{},
			}
			n.nodes.Store(name, child)
		}
		n.mu.Unlock()
	}
	child.(*node).series.add(event.Timestamp, 1)
	child.(*node).addChildNode(event, names[1:])
}

func (n *node) findNode(names []string) *node {
	if len(names) == 0 {
		return n
	}

	child, found := n.nodes.Load(names[0])
	if !found {
		return nil
	}
	return child.(*node).findNode(names[1:])
}

func (t *tree) find(query *Query) *timeSeries {
	names := sortAttributes(query.Attributes)
	node := t.root.findNode(names)
	if node == nil {
		return nil
	} else {
		return node.series
	}
}

func newTree() *tree {
	return &tree{
		root: &node{
			nodes: &sync.Map{},
		},
	}
}

func sortAttributes(attributes map[string]string) []string {
	names := make([]string, len(attributes))
	for name := range attributes {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}
