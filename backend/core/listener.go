package core

import (
	"errors"
	"spyal/contracts"
	"strconv"
)

func Listen(l contracts.Listener) error {
	echan := Hub.GetHub(l.GetEventName())
	if echan == nil {
		return errors.New("Channel " + strconv.Itoa(int(l.GetEventName())) + " is null")
	}
	go func() {
		for e := range echan {
			l.Handle(e)
		}
	}()
	return nil
}
