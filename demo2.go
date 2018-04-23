package main

import (
	"demo2/simulator"
	"fmt"
)

func main() {

/*
	world := simulator.GetWorldFromFile("4by4.map")
	path, distance := world.ShortestPath(18, 26)
	fmt.Println(path, " ", distance)
	world.SetFrameRate(10)
	world.AddCars(5)
	go world.Loop()
//*/

// TODO: remove this test code in favor of world_test.go
//*
  world := simulator.GetWorldFromFile("edgeFindTest.map")
  testcoords := simulator.Coordinates{X:120, Y:50}
  intersect, edgeID := world.ClosestEdgeAndCoord(testcoords)
  fmt.Println("Closest intersect point: ", intersect, "edgeID: ", edgeID)
//*/
}
