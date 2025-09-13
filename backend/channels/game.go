package channels

import (
	"spyal/core"
	"spyal/contracts"
)

type Game struct{
	contracts.Channel
}

func NewGameChannel() contracts.Channel {
	return &Game{
		Channel: core.NewChannel("game"),
	}
}
