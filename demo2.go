package main
import (
	"log"
	"fmt"
)


var world World
func main() {


	world := getWorldFromFile("4by4.map")
	best, err := world.graph.Shortest(18,26)
	if err!=nil{
		log.Fatal(err)
	}
	fmt.Println("Shortest distance ", best.Distance, " following path ", best.Path)
}

//// TODO
//func findClosestEdgeToPoint(randomPoint Point, edges []Edge) (Point, Edge) {
//	var closestEdge Edge
//	var closestPointOnEdge Point
//	return closestPointOnEdge, closestEdge
//}

