package bootstrap

import (
	"spyal/contracts"
	"spyal/core"
	"spyal/listeners"
)

type Initializer func() error
type NewListener func() contracts.Listener

//nolint:gochecknoglobals
var Inits []Initializer

func InitListeners() error {
	for _, newListener := range listeners.ListenerRegistry {
		listener := newListener()
		err := core.Listen(listener)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitAll() error {
	Inits = []Initializer{
		core.InitMsgHub,
		InitListeners,
	}

	for _, init := range Inits {
		err := init()
		if err != nil {
			return err
		}
	}

	return nil
}
