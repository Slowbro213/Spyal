package listeners

import (
	"encoding/json"
	"spyal/broadcasting"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"
)

type LeftListener struct {
	contracts.Listener
}

func NewLeftListener() contracts.Listener {
	return &LeftListener{}
}

func (el *LeftListener) GetEventName() contracts.EventName {
	return events.Leftevent
}

func (el *LeftListener) Handle(e contracts.Event) {
	data := e.GetData()
	err := broadcasting.Broadcast(e)
	if err != nil {
		core.Logger.Warn("LeftListener: failed to broadcast: " + err.Error())
		return
	}

	if _, err := json.Marshal(data); err != nil {
		core.Logger.Warn("LeftListener: failed to marshal data to JSON: " + err.Error())
		return
	}
}
