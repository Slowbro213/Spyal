package core

import (
	"net/http"
	"spyal/contracts"
	"strings"
	"time"

	"context"
	"io"

	"github.com/coder/websocket"

	"go.uber.org/zap"
)

const (
	timeOut = 10
)

func IsWebSocketRequest(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

type PokedServer struct {
	Log       *zap.Logger
	EmitEvent func(map[string]any) contracts.Eventer
}

func (ps *PokedServer) StartWSServer(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"http://localhost:*", "http://192.168.1.23:*"},
		Subprotocols:   []string{"poked"},
	})

	if err != nil {
		ps.Log.Error("", zap.Error(err))
		return
	}
	defer c.CloseNow() //nolint

	if c.Subprotocol() != "poked" {
		c.Close(websocket.StatusPolicyViolation, "client must speak the poked subprotocol")
		return
	}

	for {
		err = poke(c, ps.EmitEvent)
		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			ps.Log.Error("failed to echo with %v: %v"+r.RemoteAddr, zap.Error(err))
			return
		}
	}
}

func poke(c *websocket.Conn, emitEvent func(map[string]any) contracts.Eventer) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeOut)
	defer cancel()

	typ, r, err := c.Reader(ctx)
	if err != nil {
		return err
	}

	// Read all of `r` into memory
	payload, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	// Build the event using the raw message
	data := map[string]any{
		"msg": string(payload),
	}

	event := emitEvent(data)
	go Dispatch(event)

	// Echo the same message back to the client
	w, err := c.Writer(ctx, typ)
	if err != nil {
		return err
	}

	_, err = w.Write(payload)
	if err != nil {
		return err
	}

	return w.Close()
}
