package listeners

import (
	"encoding/json"
	"spyal/broadcasting"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"
)

type GameEndListener struct {
	contracts.Listener
}

func NewGameEndListener() contracts.Listener {
	return &GameEndListener{}
}

func (el *GameEndListener) GetEventName() contracts.EventName {
	return events.Gameendevent
}

func (el *GameEndListener) Handle(e contracts.Event) {
	data := e.GetData()
	err := broadcasting.Broadcast(e)
	if err != nil {
		core.Logger.Warn("GameEndListener: failed to broadcast: " + err.Error())
		return
	}

	if _, err := json.Marshal(data); err != nil {
		core.Logger.Warn("GameEndListener: failed to marshal data to JSON: " + err.Error())
		return
	}
}
