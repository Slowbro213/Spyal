package listeners

import (
	"encoding/json"
	"spyal/broadcasting"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"
)

type EchoListener struct {
	contracts.Listener
}

func NewEchoListener() contracts.Listener {
	return &EchoListener{}
}

func (el *EchoListener) GetEventName() contracts.EventName {
	return events.Echoevent
}

func (el *EchoListener) Handle(e contracts.Event) {
	data := e.GetData()
	err := broadcasting.Broadcast(e)
	if err != nil {
		core.Logger.Warn("EchoListener: failed to broadcast: " + err.Error())
		return
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		core.Logger.Warn("EchoListener: failed to marshal data to JSON: " + err.Error())
		return
	}
	core.Logger.Info("New Data JSON: " + string(jsonBytes))
}
