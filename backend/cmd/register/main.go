package main

import (
	"log"
	"spyal/tooling"
)

const (
	listenersDir = "listeners"
	eventsDir = "events"
	channelsDir = "channels"
)

//gosec:disable G101 -- This is a false positive
func main() {
	dirs := [1]string{listenersDir}

	for _, dir := range dirs {
		constructors, err := tooling.GetConstructors(dir)

		if err != nil {
			log.Fatalf("Error getting constructors %v", err)
		}

		err = tooling.GenerateConstructorRegistry(dir, constructors)
		if err != nil {
			log.Fatalf("Error generating constructor registry %v", err)
		}
	}

	events, err := tooling.GetEventTypes(eventsDir)

	if err != nil {
		log.Fatalf("Error while getting events %v", err)
	}

	err = tooling.GenerateEventTypes(eventsDir, events)

	if err != nil {
		log.Fatalf("Error while generating events registry %v", err)
	}


	channels, err := tooling.GetChannelConstructors(channelsDir)

	if err != nil {
		log.Fatalf("Error while getting channels %v", err)
	}

	err = tooling.GenerateChannelRegistryFile(channelsDir,channels)

	if err != nil {
		log.Fatalf("Error while generating channels registry %v", err)
	}
}
