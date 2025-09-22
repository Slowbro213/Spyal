package events

import (
	"spyal/contracts"
)


type ChatEvent struct {
	BaseEvent
}

func NewChatEvent(data map[string]any) contracts.Event{
	return &ChatEvent{
		BaseEvent: BaseEvent{
			Data: data,
		},
	}
}


func (ee *ChatEvent) GetName() contracts.EventName {
	return Chatevent
}

func (ee *ChatEvent) Channel() string{
	return "game"
}

func (ee *ChatEvent) Topic() string {
	v, _ := ee.Data["topic"].(string)
	return v
}
