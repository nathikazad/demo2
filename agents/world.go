package agents

import (
	"demo2/dijkstra"
	"time"
	"math/rand"
)

type Edge struct {
	id uint
	startingNode *Node
	endingNode   *Node
	weight int
}

type Node struct {
	id uint
	coordinates Coordinates
	nextNodes []*Node
	nextEdges map[*Node]*Edge
}

type World struct {
	graph *dijkstra.Graph
	nodes map[uint]*Node
	edges map[uint]*Edge
	cars map[uint]*Car
	carStates map[*Car]CarState
	carCommands map[*Car]CarCommand
	broadcast chan map[*Car]CarState
	commandReceiver chan CarCommand
	frameRate uint
	nextCarIndex uint
}

func NewWorld() *World {
	world := new(World)
	world.frameRate = 1
	world.nextCarIndex = 0
	world.cars = make(map[uint]*Car)
	world.carStates = make(map[*Car]CarState)
	world.carCommands = make(map[*Car]CarCommand)
	return world
}

// TODO: toss this for our own concurrency safe implementation
func (w *World) ShortestPath(n1_id uint, n2_id uint) ([]int, uint){
	best, _ := w.graph.Shortest(int(n1_id), int(n2_id))
	return best.Path, uint(best.Distance)
}

// ClosestEdgeAndCoord For coords within world space, find  closest coords on an edge on world graph
func (w *World) ClosestEdgeAndCoord(queryPoint Coordinates) (intersect Coordinates, edgeID uint) {
	// TODO: implement me
	return
}

func (w *World) SetFrameRate(frameRate uint) {
	w.frameRate = frameRate
}

func (w World) AddCars(numberOfCars int) ([]*Car) {
	cars := make([]*Car, numberOfCars)
	for i := 0; i < numberOfCars; i++{
		cars = append(cars, w.AddCar())
	}
	return cars
}

func (w World) AddCar() (*Car) {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	edge := w.edges[uint(r1.Int() % len(w.edges))]
	car := Car{id:w.nextCarIndex, path:[]uint{} }
	w.cars[car.id] = &car
	//TODO check for collison and add orientation
	w.carStates[&car] =  CarState{id:car.id, coordinates:edge.startingNode.coordinates, orientation:0, edgeId:edge.id, action:Stopped}
	w.nextCarIndex++
	return &car
}

func (w World) Loop() {
	for _, car := range w.cars {
		car.startLoop(w.broadcast, w.commandReceiver)
	}
	for {
		w.broadcast <- w.carStates
		timeToWait := time.Duration(1000/w.frameRate) * time.Millisecond
		timer := time.NewTimer(timeToWait)
		for {
			select {
			case command := <-w.commandReceiver:
				w.carCommands[w.cars[command.id]] = command
				continue
			case <-timer.C:
				break
			}
		}
		w.executeCommands()
	}
}

func (w World) executeCommands() {

}
