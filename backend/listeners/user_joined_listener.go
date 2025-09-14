package listeners

import (
	"encoding/json"
	"spyal/broadcasting"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"
)

type UserJoinedListener struct {
	contracts.Listener
}

func NewUserJoinedListener() contracts.Listener {
	return &UserJoinedListener{}
}

func (el *UserJoinedListener) GetEventName() contracts.EventName {
	return events.Userjoinedevent
}

func (el *UserJoinedListener) Handle(e contracts.Event) {
	data := e.GetData()
	err := broadcasting.Broadcast(e)
	if err != nil {
		core.Logger.Warn("UserJoinedListener: failed to broadcast: " + err.Error())
		return
	}
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		core.Logger.Warn("UserJoinedListener: failed to marshal data to JSON: " + err.Error())
		return
	}
	core.Logger.Info("New Data JSON: " + string(jsonBytes))
}
