package listeners

import (
	"encoding/json"
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

func (el *EchoListener) Handle(data map[string]any) {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		core.Logger.Warn("EchoListener: failed to marshal data to JSON: " + err.Error())
		return
	}
	core.Logger.Info("New Data JSON: " + string(jsonBytes))
}
