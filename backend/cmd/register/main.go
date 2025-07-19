package main

import (
	"log"
	"spyal/tooling"
)

func main() {
	dirs := [1]string{"listeners"}

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

	events, err := tooling.GetEventTypes("events")

	if err != nil {
		log.Fatalf("Error while getting events %v", err)
	}

	err = tooling.GenerateEventTypes("events", events)

	if err != nil {
		log.Fatalf("Error while generating events %v", err)
	}
}
