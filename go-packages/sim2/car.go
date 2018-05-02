package sim2

import (
  //"fmt"
  "math"
)

// car - Describes routine hooks and logic for Cars within a World simulation

// Car - struct for all info needed to manage a Car within a World simulation.
type Car struct {
  // TODO determine if Car needs any additional/public members

  id uint
  dstNode uint
  path []uint
  world *World
  syncChan chan bool
  sendChan chan CarInfo
}

// NewCar - Construct a new valid Car object
func NewCar(id uint, w *World, sync chan bool, send chan CarInfo) *Car {
  c := new(Car)
  c.id = id
  c.world = w
  c.syncChan = sync
  c.sendChan = send
  return c
}

// CarLoop - Begin the car simulation execution loop
func (c *Car) CarLoop() {
  for {
    // Block waiting for next sync event
    <-c.syncChan
    //fmt.Println("Car", c.id, ": got sync")

    // Check for new solidity bids and update the bid queue
    // TODO: this

    // Evaluate destination based on current location and status of bids
    // TODO: this

    // Check for shortest path to destination
    // TODO: this

    // Evaluate road rules and desired movement within world
    // TODO: this

    // Send movement update request to World
    // TODO: replace this with real update
    c.sendChan <- CarInfo{ Pos:Coords{0,0}, Vel:Coords{0,0} }
    //fmt.Println("Car", c.id, ": sent update")
  }
}

// closestEdgeAndCoord For coords within world space, find  closest coords on an edge on world graph
// Return coordinates of closest point on world graph, and corresponding edge ID in world struct
func (w *World) closestEdgeAndCoord(queryPoint Coords) (intersect Coords, edgeID uint) {
  // TODO: input sanitation/validation; error handling?
  // TODO: proper helper function breakdown of closestEdgeAndCoord

  shortestDistance := math.Inf(1)
  intersect = Coords{0, 0}
  edgeID = 0

  // TODO: remove randomness caused by traversing equivalent closest edges with 'range' on map here
  for edgeIdx, edge := range w.Graph.Edges {
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
func (e *Edge) checkIntersect(query Coords) (intersect Coords, distance float64) {
  // Via https://stackoverflow.com/a/5205747

  // TODO: cleaner logic in implementation, needs to be several helper functions for testing

  // Query coords as floats for internal math
  xQuery := query.X
  yQuery := query.Y

  // Edge coordinate references
  x1 := e.Start.Pos.X
  y1 := e.Start.Pos.Y
  x2 := e.End.Pos.X
  y2 := e.End.Pos.Y

  // Slope of edge and perpendecular section
  mEdge := (y2 - y1) / (x2 - x1)
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
    distance1 := query.Distance(e.Start.Pos)
    distance2 := query.Distance(e.End.Pos)
    //fmt.Println(", dist1: ", distance1, ", dist2: ", distance2)
    if (distance1 < distance2) {
      intersect = e.Start.Pos
      distance = distance1
    } else {
      intersect = e.End.Pos
      distance = distance2
    }
  } else {  // In perpendicular region; find intersection on Edge
    // Check for straight lines to ease calculations
    if (x1 == x2) {
      intersect.X = math.Round(x1)
      intersect.Y = math.Round(yQuery)

    } else if (y1 == y2) {
      intersect.X = math.Round(xQuery)
      intersect.Y = math.Round(y1)
    } else {
      intX := (mEdge * x1 - mPerp * xQuery + yQuery - y1) / (mEdge - mPerp)
      intY := mPerp * (intersect.X - xQuery) + yQuery
      intersect.X = math.Round(intX)
      intersect.Y = math.Round(intY)
    }
    distance = intersect.Distance(query)
  }
  //fmt.Print(", intersect: ", intersect, ", distance: ", distance)
  //fmt.Println()
  return
}
