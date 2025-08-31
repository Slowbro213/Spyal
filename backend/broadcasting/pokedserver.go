package broadcasting

import (
	"net/http"
	"os"
	"fmt"
	"errors"
	"encoding/json"
	"spyal/channels"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"
	"strings"
	"time"

	"context"

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
	EmitEvent func(map[string]any) contracts.Event
}



func (ps *PokedServer) StartWSServer(w http.ResponseWriter, r *http.Request) {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: []string{"https?://" + host + ":" + port},
		Subprotocols:   []string{"poked"},
	})
	if err != nil {
		ps.Log.Error("failed to accept websocket", zap.Error(err))
		return
	}

	// defer CloseNow and check error
	defer func() {
		if err := c.CloseNow(); err != nil {
			ps.Log.Warn("failed to close websocket", zap.Error(err))
		}
	}()

	conn := core.NewWSConnection(c, r.RemoteAddr)

	if conn.Subprotocol() != "poked" {
		if err := conn.Close(websocket.StatusPolicyViolation, "client must speak the poked subprotocol"); err != nil {
			ps.Log.Warn("failed to close websocket on subprotocol mismatch", zap.Error(err))
		}
		return
	}

	channelName := r.URL.Query().Get("channel")
	chann := channels.Channels[channelName]

	if chann == nil {
		return
	}

	if ok := chann.Join(conn, r); !ok {
		ps.Log.Warn("failed to join channel", zap.Error(err))
	} 

	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*timeOut)
		err := poke(ctx, conn)
		cancel()

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure {
			return
		}
		if err != nil {
			ps.Log.Error("connection error with "+conn.RemoteAddr(), zap.Error(err))
			return
		}
	}
}



func poke(ctx context.Context, conn contracts.WSConnection) error {
	// Read the payload from the connection
	payload, err := conn.Read(ctx)
	if err != nil {
		return err
	}

	// Decode JSON into a generic map
	var msg map[string]any
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	// Extract "type" field and cast to EventName
	rawType, ok := msg["type"]
	if !ok {
		return errors.New("missing type field in message")
	}

	var eventType contracts.EventName
	switch v := rawType.(type) {
	case float64:
		eventType = contracts.EventName(int(v))
	case int:
		eventType = contracts.EventName(v)
	default:
		return fmt.Errorf("invalid type field: %T", rawType)
	}

	event := events.NewEvent(eventType, msg)
	if event == nil {
		return fmt.Errorf("unknown event type: %d", eventType)
	}

	go core.Dispatch(event)
	return err
}

