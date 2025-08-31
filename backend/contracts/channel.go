package contracts

import "net/http"

type Channel interface {
	Name() string
	Join(WSConnection, *http.Request) bool
	Leave(WSConnection) bool
	WSConnections() map[int]WSConnection
}
