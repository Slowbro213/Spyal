package listeners

import (
	"encoding/json"
	"spyal/broadcasting"
	"spyal/contracts"
	"spyal/core"
	"spyal/events"
)

type ChatListener struct {
	contracts.Listener
}

func NewChatListener() contracts.Listener {
	return &ChatListener{}
}

func (el *ChatListener) GetEventName() contracts.EventName {
	return events.Chatevent
}

func (el *ChatListener) Handle(e contracts.Event) {
	data := e.GetData()
	err := broadcasting.Broadcast(e)
	if err != nil {
		core.Logger.Warn("ChatListener: failed to broadcast: " + err.Error())
		return
	}
	if _, err := json.Marshal(data); err != nil {
		core.Logger.Warn("ChatListener: failed to marshal data to JSON: " + err.Error())
		return
	}
}
