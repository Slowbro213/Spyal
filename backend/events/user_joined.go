package events

import (
	"spyal/contracts"
)

type UserJoinedEvent struct {
	BaseEvent
}

func NewUserJoinedEvent(data map[string]any) contracts.Event {
	eventData := map[string]any{
		"type": Userjoinedevent,
		"msg":  data,
	}
	return &UserJoinedEvent{
		BaseEvent: BaseEvent{
			Data: eventData,
		},
	}
}

func (ee *UserJoinedEvent) GetName() contracts.EventName {
	return Gameevent
}

func (ee *UserJoinedEvent) Channel() string {
	return "game"
}

func (ee *UserJoinedEvent) Topic() string {
	v, _ := ee.Data["topic"].(string)
	return v
}
