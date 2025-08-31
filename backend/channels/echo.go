package channels

import (
	"spyal/core"
	"spyal/contracts"
)

type Echo struct{
	contracts.Channel
}

func NewEchoChannel() contracts.Channel {
	return &Echo{
		Channel: core.NewChannel("echo"),
	}
}
