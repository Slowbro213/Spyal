package broadcasting

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"spyal/channels"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"

	"github.com/coder/websocket"
	"go.uber.org/zap"
)

const (
	timeOutSeconds = 10_000_000
)

func IsWebSocketRequest(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

type PokedServer struct {
	Log       *zap.Logger
	EmitEvent func(map[string]any) contracts.Event
}



//nolint
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
	topic := r.URL.Query().Get("topic")
	chann := channels.Channels[channelName]

	if chann == nil {
		ps.Log.Warn("unknown channel", zap.String("channel", channelName))
		return
	}

	if ok := chann.Join(conn, topic, r); !ok {
		ps.Log.Warn("failed to join channel")
	}

	for {
		ctx, cancel := context.WithTimeout(r.Context(), time.Second*time.Duration(timeOutSeconds))
		err := poke(ctx, conn, topic)
		cancel()

		if err == nil {
			continue
		}

		if websocket.CloseStatus(err) == websocket.StatusNormalClosure ||
			websocket.CloseStatus(err) == websocket.StatusGoingAway ||
			websocket.CloseStatus(err) == websocket.StatusNoStatusRcvd {
			ps.Log.Info("websocket closed by client",
				zap.String("remote", conn.RemoteAddr()),
				zap.Int("close_code", int(websocket.CloseStatus(err))),
			)
			chann.Leave(conn, topic)
			return
		}

		if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) || errors.Is(err, io.EOF) ||
			strings.Contains(err.Error(), "use of closed network connection") {
			ps.Log.Info("websocket read ended (client disconnected / context done)",
				zap.String("remote", conn.RemoteAddr()),
				zap.Error(err),
			)
			chann.Leave(conn, topic)
			return
		}

		ps.Log.Error("connection error with "+conn.RemoteAddr(), zap.Error(err))
		chann.Leave(conn, topic)
		return
	}
}

func poke(ctx context.Context, conn contracts.WSConnection, topic string) error {
	payload, err := conn.Read(ctx)
	if err != nil {
		return err
	}

	var msg map[string]any
	if err := json.Unmarshal(payload, &msg); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	msg["topic"] = topic

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

	core.Dispatch(event)
	return nil
}
