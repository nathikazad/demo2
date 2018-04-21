package main

import (
	"demo2/agents"
)

func main() {

	world := agents.GetWorldFromFile("4by4.map")
	world.SetFrameRate(10)
	world.AddCars(5)
	go world.Loop()
}


