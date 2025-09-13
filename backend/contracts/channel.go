package contracts

import "net/http"

type Channel interface {
	Name() string
	Join(WSConnection,string, *http.Request) bool
	Leave(WSConnection,string) bool
	WSConnections() map[int]WSConnection
}
