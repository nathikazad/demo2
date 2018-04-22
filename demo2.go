package main

import (
	"demo2/simulator"
	"fmt"
)

func main() {

	world := simulator.GetWorldFromFile("4by4.map")
	path, distance := world.ShortestPath(18, 26)
	fmt.Println(path, " ", distance)
	world.SetFrameRate(10)
	world.AddCars(5)
	go world.Loop()
}


