package simulator

import (
	"demo2/dijkstra"
	"time"
	"math/rand"
  "math"
  //"fmt"
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

// TODO: remove commented-out test prints
// TODO: proper helper function breakdown of ClosestEdgeAndCoord
// TODO: proper test file world_test.go

// ClosestEdgeAndCoord For coords within world space, find  closest coords on an edge on world graph
// Return coordinates of closest point on world graph, and corresponding edge ID in world struct
func (w *World) ClosestEdgeAndCoord(queryPoint Coordinates) (intersect Coordinates, edgeID uint) {
  // TODO: input sanitation/validation; error handling?

  shortestDistance := math.Inf(1)
  intersect = Coordinates{0, 0}
  edgeID = 0

  // TODO: remove randomness caused by traversing equivalent closest edges with 'range' on map here
  for edgeIdx, edge := range w.edges {
    coord, dist := edge.checkIntersect(queryPoint)
    //fmt.Print("[", edgeIdx, "] <ClosestEdgeAndCoord>")
    //fmt.Print(" query: ", queryPoint, ", shortest: ", shortestDistance, ", new: ", coord, ", dist: ", dist)
    if dist < shortestDistance {
      shortestDistance = dist
      intersect = coord
      edgeID = edgeIdx
      //fmt.Print(" (new shortest: ", dist , " @", coord, ")")
    }
    //fmt.Println()
  }
  return
}

// checkIntersect find the point on this edge closest to the query coordinates and report distance
// Return coordinates of closest point and distance of closest point to query
func (e *Edge) checkIntersect(query Coordinates) (intersect Coordinates, distance float64) {
  // Via https://stackoverflow.com/a/5205747

  // TODO: cleaner logic in implementation, needs to be several helper functions for testing

  // Query coords as floats for internal math
  xQuery := float64(query.X)
  yQuery := float64(query.Y)

  // Edge coordinate references
  x1 := float64(e.startingNode.coordinates.X)
  y1 := float64(e.startingNode.coordinates.Y)
  x2 := float64(e.endingNode.coordinates.X)
  y2 := float64(e.endingNode.coordinates.Y)

  // Slope of edge and perpendecular section
  mEdge := float64(y2 - y1) / float64(x2 - x1)
  mPerp := float64(-1) / mEdge

  //fmt.Print("<checkIntersect>")
  //fmt.Print(" query: {", xQuery, ",", yQuery, "}")
  //fmt.Print(", Edge{", x1, ",", y1, "}->{", x2, ",", y2, "}")
  //fmt.Print(", mEdge: ", mEdge, ", mPerp: ", mPerp)

  inPerp := false  // By default, not in perpendicular region
  isVert := math.IsInf(mPerp, 0)  // Is the perpendicular section vertical?

  // Determine if the query point lies in region perpendicular to the edge segment
  if (isVert) {  // Regions determined by x values
    if ((x1 < x2 && x1 < xQuery && xQuery < x2) || (x2 < x1 && xQuery < x1 && x2 < xQuery)) {
      inPerp = true
    }
  } else {  // Regions determined by y values
    // Relative y-coordinates of perpendicular lines at bounds of edge segment at x-value of query
    y1Rel := mPerp * (xQuery - x1) + y1
    y2Rel := mPerp * (xQuery - x2) + y2
    if ((y1Rel < y2Rel && y1Rel < yQuery && yQuery < y2Rel) || (y2Rel < y1Rel && yQuery < y1Rel && y2Rel < yQuery)) {
      inPerp = true
    }
  }

  //fmt.Print(", inPerp: ", inPerp, ", isVert: ", isVert)

  if (!inPerp) {  // Not in perpendicular segment or vertical; identify closest
    distance1 := query.Distance(e.startingNode.coordinates)
    distance2 := query.Distance(e.endingNode.coordinates)
    //fmt.Println(", dist1: ", distance1, ", dist2: ", distance2)
    if (distance1 < distance2) {
      intersect = e.startingNode.coordinates
      distance = distance1
    } else {
      intersect = e.endingNode.coordinates
      distance = distance2
    }
  } else {  // In perpendicular region; find intersection on Edge
    // Check for straight lines to ease calculations
    if (x1 == x2) {
      intersect.X = uint(x1)
      intersect.Y = uint(yQuery)

    } else if (y1 == y2) {
      intersect.X = uint(xQuery)
      intersect.Y = uint(y1)
    } else {
      intersect.X = uint((mEdge * x1 - mPerp * xQuery + yQuery - y1) / (mEdge - mPerp))
      intersect.Y = uint(mPerp * (float64(intersect.X) - xQuery) + yQuery)
    }
    distance = intersect.Distance(query)
  }
  //fmt.Print(", intersect: ", intersect, ", distance: ", distance)
  //fmt.Println()
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
