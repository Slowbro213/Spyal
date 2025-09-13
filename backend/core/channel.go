package core

import (
	"net/http"
	"spyal/contracts"
	"sync"
	"maps"
)

type topic struct {
	wsConnections map[int]contracts.WSConnection
	count         int
	mu            sync.RWMutex
}

type Channel struct {
	name   string
	topics map[string]*topic
	mu     sync.RWMutex
}

func NewChannel(name string) contracts.Channel {
	return &Channel{
		name:   name,
		topics: make(map[string]*topic),
	}
}

func (c *Channel) Name() string { return c.name }

func (c *Channel) Join(wsc contracts.WSConnection, top string, r *http.Request) bool {
	if !c.auth(r) {
		return false
	}

	c.mu.Lock()
	t, ok := c.topics[top]
	if !ok {
		t = &topic{
			wsConnections: make(map[int]contracts.WSConnection),
			count:         0,
		}
		c.topics[top] = t
	}
	c.mu.Unlock()

	t.mu.Lock()
	id := t.count
	t.wsConnections[id] = wsc
	t.count++
	t.mu.Unlock()

	return true
}

func (c *Channel) Leave(wsc contracts.WSConnection, top string) bool {
	c.mu.RLock()
	t, ok := c.topics[top]
	c.mu.RUnlock()
	if !ok {
		return false
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	for id, conn := range t.wsConnections {
		if conn == wsc {
			delete(t.wsConnections, id)
			return true
		}
	}
	return false
}

func (c *Channel) WSConnections() map[int]contracts.WSConnection {
	c.mu.RLock()
	defer c.mu.RUnlock()

	out := make(map[int]contracts.WSConnection)
	for _, t := range c.topics {
		t.mu.RLock()
		maps.Copy(out, t.wsConnections)
		t.mu.RUnlock()
	}
	return out
}

func (c *Channel) auth(_ *http.Request) bool {
	return true
}

