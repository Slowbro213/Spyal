package core

import (
	"maps"
	"net/http"
	"spyal/contracts"
	"sync"
)

type Channel struct {
	name          string
	count         int
	wsConnections map[int]contracts.WSConnection
	mu            sync.RWMutex
}

func NewChannel(name string) contracts.Channel {
	return &Channel{
		name:          name,
		count:         0,
		wsConnections: make(map[int]contracts.WSConnection),
	}
}

func (c *Channel) Name() string {
	return c.name
}

func (c *Channel) Join(wsc contracts.WSConnection, r *http.Request) bool {
	authenticated := c.auth(r)

	if !authenticated {
		return false
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.wsConnections[c.count] = wsc
	c.count++
	return true
}

func (c *Channel) Leave(wsc contracts.WSConnection) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	for id, conn := range c.wsConnections {
		if conn == wsc {
			delete(c.wsConnections, id)
			return true
		}
	}

	return false
}

func (c *Channel) WSConnections() map[int]contracts.WSConnection {
	c.mu.RLock()
	defer c.mu.RUnlock()

	m := make(map[int]contracts.WSConnection, len(c.wsConnections))
	maps.Copy(m, c.wsConnections)
	return m
}

func (c *Channel) auth(_ *http.Request) bool {
	return true
}
