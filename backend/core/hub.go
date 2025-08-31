package core

import (
	"spyal/contracts"
	"spyal/events"
)

type MsgHub struct {
	hubs []chan map[string]any
}

func (mh *MsgHub) GetHub(e contracts.EventName) chan map[string]any {
	return mh.hubs[e]
}

//nolint:gochecknoglobals
var Hub *MsgHub

func InitMsgHub() error {
	if Hub != nil {
		return nil
	}

	Hub = &MsgHub{
		hubs: make([]chan map[string]any, int(events.EventNameCount)),
	}

	for i := range int(events.EventNameCount) {
		Hub.hubs[i] = make(chan map[string]any)
	}

	return nil
}

func Dispatch(e contracts.Event) {
	echan := Hub.GetHub(e.GetName())

	echan <- e.GetData()
}
