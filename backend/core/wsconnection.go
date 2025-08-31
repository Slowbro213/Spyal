package core

import (
	"context"
	"io"
	"spyal/contracts"

	"github.com/coder/websocket"
	"go.uber.org/zap"
)

type WSConnection struct {
	c *websocket.Conn
	remoteAddr string
}

func NewWSConnection(c *websocket.Conn, remoteAddr string) contracts.WSConnection {
	return &WSConnection{c: c, remoteAddr: remoteAddr}
}

func (ws *WSConnection) Read(ctx context.Context) ([]byte, error) {
	_, r, err := ws.c.Reader(ctx)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(r)
}

func (ws *WSConnection) Write(ctx context.Context, payload []byte) error {
	w, err := ws.c.Writer(ctx, websocket.MessageText)
	if err != nil {
		return err
	}
	_, err = w.Write(payload)
	if err != nil {
		return err
	}
	return w.Close()
}

func (ws *WSConnection) Close(code websocket.StatusCode, reason string) error {
	return ws.c.Close(code, reason)
}

func (ws *WSConnection) CloseNow() {
	err := ws.c.CloseNow()
	if err != nil {
		Logger.Error("Error closing connection from " + ws.remoteAddr, zap.Error(err))
	}
}

func (ws *WSConnection) Subprotocol() string {
	return ws.c.Subprotocol()
}

func (ws *WSConnection) RemoteAddr() string {
	return ws.remoteAddr
}
