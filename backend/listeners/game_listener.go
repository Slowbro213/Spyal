package listeners

import (
	"encoding/json"
	"spyal/broadcasting"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"
)

type GameListener struct {
	contracts.Listener
}

func NewGameListener() contracts.Listener {
	return &GameListener{}
}

func (el *GameListener) GetEventName() contracts.EventName {
	return events.Gameevent
}

func (el *GameListener) Handle(e contracts.Event) {
	data := e.GetData()
	err := broadcasting.Broadcast(e)
	if err != nil {
		core.Logger.Warn("GameListener: failed to broadcast: " + err.Error())
		return
	}

	if _, err := json.Marshal(data); err != nil {
		core.Logger.Warn("GameListener: failed to marshal data to JSON: " + err.Error())
		return
	}
}
