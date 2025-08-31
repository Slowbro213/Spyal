package core

import (
	"spyal/contracts"
	"spyal/events"
)

type MsgHub struct {
	hubs []chan contracts.Event
}

func (mh *MsgHub) GetHub(e contracts.EventName) chan contracts.Event {
	return mh.hubs[e]
}

//nolint:gochecknoglobals
var Hub *MsgHub

func InitMsgHub() error {
	if Hub != nil {
		return nil
	}

	Hub = &MsgHub{
		hubs: make([]chan contracts.Event, int(events.EventNameCount)),
	}

	for i := range int(events.EventNameCount) {
		Hub.hubs[i] = make(chan contracts.Event)
	}

	return nil
}

func Dispatch(e contracts.Event) {
	echan := Hub.GetHub(e.GetName())

	echan <- e
}
