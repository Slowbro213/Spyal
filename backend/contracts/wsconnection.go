package contracts

import (
	"context"

	"github.com/coder/websocket"
)


type WSConnection interface {
	Read(ctx context.Context) ([]byte, error)
	Write(ctx context.Context, payload []byte) error
	Close(code websocket.StatusCode, reason string) error
	CloseNow()
	Subprotocol() string
	RemoteAddr() string
}
