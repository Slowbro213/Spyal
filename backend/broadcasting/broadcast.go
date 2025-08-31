package broadcasting

import (
	"context"
	"encoding/json"
	"errors"
	"spyal/channels"
	"spyal/contracts"
	"spyal/core"
	"time"

	"go.uber.org/zap"
)

const broadcastTimeout = 10 // seconds

func Broadcast(e contracts.Event) error {
	// If the event implements ShouldBroadcast, call Broadcast
	b, ok := e.(contracts.ShouldBroadcast)
	if !ok {
		return errors.New("event not broadcastable")
	}

	channelName := b.Channel()
	data := e.GetData()

	chann := channels.Channels[channelName]
	conns := chann.WSConnections()

	// Prepare the payload to send (for simplicity, just JSON)
	payload, err := json.Marshal(data)
	if err != nil {
	    return err
	}

	// Iterate over all connections
	for _, wsc := range conns {
		// Wrap with context for timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*broadcastTimeout)
		err := wsc.Write(ctx, payload)
		cancel()
		if err != nil {
			core.Logger.Error("Error broadcasting event ", zap.Error(err))
		}
	}

	return nil
}
