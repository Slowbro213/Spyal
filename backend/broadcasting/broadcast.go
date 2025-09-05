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

const broadcastTimeout = 10

func Broadcast(e contracts.Event) error {
	b, ok := e.(contracts.ShouldBroadcast)
	if !ok {
		return errors.New("event not broadcastable")
	}

	channelName := b.Channel()
	data := e.GetData()

	chann := channels.Channels[channelName]
	conns := chann.WSConnections()

	payload, err := json.Marshal(data)
	if err != nil {
	    return err
	}

	for _, wsc := range conns {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*broadcastTimeout)
		err := wsc.Write(ctx, payload)
		cancel()
		if err != nil {
			core.Logger.Error("Error broadcasting event ", zap.Error(err))
		}
	}

	return nil
}
