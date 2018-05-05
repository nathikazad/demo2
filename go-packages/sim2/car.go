package sim2

import (
  //"fmt"
  "math"
  "github.com/ethereum/go-ethereum/ethclient"
  "github.com/ethereum/go-ethereum/accounts/abi/bind"
  "github.com/ethereum/go-ethereum/common"
  "log"
  "github.com/ethereum/go-ethereum/crypto"
  "github.com/ethereum/go-ethereum/core/types"
  "context"
  "math/rand"
  "time"
  "fmt"
)

// car - Describes routine hooks and logic for Cars within a World simulation

// Car - struct for all info needed to manage a Car within a World simulation.
type Car struct {
  // TODO determine if Car needs any additional/public members
  id uint
  path Path
  world *World
  syncChan chan bool
  sendChan chan CarInfo
  ethApi EthAPI
  acceptanceState AcceptanceState
}

type Path struct {
  pickUp Coords
  dropOff Coords
  currentPos Coords
  currentDir Coords
  currentEdgeId uint
  pathState PathState
  routeNodes []uint
}

type PathState int
const (
  AtRandom PathState = 0
  ToPickUp PathState = 1
  ToDropOff PathState = 2
)

type AcceptanceState int
const (
  None AcceptanceState = 0
  Trying AcceptanceState = 1
  Success AcceptanceState = 2
  Fail AcceptanceState = 3
)

type EthAPI struct {
  conn *ethclient.Client
  mrm *MoovRideManager
  auth *bind.TransactOpts
  lastEventIndex uint
  chosenRider common.Address
}

// NewCar - Construct a new valid Car object
func NewCar(id uint, w *World, sync chan bool, send chan CarInfo, mrmAddress string, privateKeyString string) *Car {
  c := new(Car)
  c.id = id
  c.world = w
  c.syncChan = sync
  c.sendChan = send
  c.acceptanceState = None


  s1 := rand.NewSource(time.Now().UnixNano())
  r1 := rand.New(s1)
  edge := c.world.Graph.Edges[uint(r1.Int() % len(c.world.Graph.Edges))]

  c.path.pathState = AtRandom
  c.path.currentPos = edge.Start.Pos
  c.path.currentDir = edge.Start.Pos.UnitVector(edge.End.Pos)
  c.path.currentEdgeId = edge.ID
  //TODO check for collison and add orientation


  var err error
  c.ethApi.conn, err = ethclient.Dial("ws://127.0.0.1:8546")
  if err != nil {
    log.Fatalf("could not create ipc client: %v", err)
  }
  c.ethApi.mrm, err = NewMoovRideManager(common.HexToAddress(mrmAddress), c.ethApi.conn)
  if err != nil {
    log.Fatalf("could not connect to mrm: %v", err)
  }

  privateKey, err := crypto.HexToECDSA(privateKeyString)
  if err != nil {
    log.Fatalf("could not convert private key to hex: %v", err)
  }
  c.ethApi.auth = bind.NewKeyedTransactor(privateKey)
  return c
}

// CarLoop - Begin the car simulation execution loop
func (c *Car) CarLoop() {

  posx := float64(0)  // TODO remove
  for {
      <-c.syncChan // Block waiting for next sync event
      //TODO: first check that car state is not in the middle of a ride
      if c.path.pathState == AtRandom {
        if c.acceptanceState != Trying {
          if c.acceptanceState == Success {
            locations, err := c.ethApi.mrm.GetLocations(nil, c.ethApi.chosenRider)
            if err != nil {
              log.Println("get locations error: %v", err)
            }
            x, y := splitCSV(locations.From)
            c.path.pickUp = Coords{x,y}
            x, y = splitCSV(locations.To)
            c.path.dropOff = Coords{x,y}
            closestCoords, closestEdgeId := c.world.closestEdgeAndCoord(c.path.pickUp)
            // Overwrite the pickup destination with the closest point on edge
            c.path.pickUp = closestCoords;
            c.path.routeNodes, _ = c.world.Graph.ShortestPath(c.world.Graph.Edges[c.path.currentEdgeId].End.ID, c.world.Graph.Edges[closestEdgeId].Start.ID)
            c.path.pathState = ToPickUp
            c.acceptanceState = None
          } else {
            addresses, err := c.ethApi.mrm.GetAvailableRides(nil)
            if err != nil {
              log.Println("could not get addresses from car: %v", err)
            }
            if len(addresses) > 0 {
              fmt.Println("Found a Guy")
              go c.tryToAcceptRequest(addresses[0])
              c.acceptanceState = Trying;
            }
          }
        }
      }

      if c.path.pathState == ToPickUp || c.path.pathState == ToDropOff {

      } else if c.path.pathState == AtRandom {
        endVertex := c.world.Graph.Edges[c.path.currentEdgeId].End
        if c.path.currentPos.Equals(endVertex.Pos) {  // Already at edge end, so change edge
          if len(endVertex.AdjEdges) == 0 {fmt.Println(endVertex.ID)}
          c.path.currentEdgeId = endVertex.AdjEdges[0].ID
          c.path.currentDir = endVertex.AdjEdges[0].Start.Pos.UnitVector(endVertex.AdjEdges[0].End.Pos)
        }else if c.path.currentPos.Distance(endVertex.Pos) > 1.0 {
          c.path.currentPos = c.path.currentPos.ProjectInDirection(1.0, endVertex.Pos)
        } else {
          c.path.currentPos = endVertex.Pos
        }
      }
      // TODO: this

      // Evaluate destination based on current location and status of bids
      // TODO: this

      // Check for shortest path to destination
      // TODO: this

      // Evaluate road rules and desired movement within world
      // TODO: this

      // Send movement update request to World
      // TODO: replace this with real update
      posx += 5
      inf := CarInfo{ Pos:c.path.currentPos, Vel:Coords{0,0}, Dir:c.path.currentDir }
      c.sendChan <- inf
      //fmt.Println("Car", c.id, ": sent update", inf)
  }
}


func (c *Car) tryToAcceptRequest (address common.Address) {
  transaction, err := c.ethApi.mrm.AcceptRideRequest(&bind.TransactOpts{
    From:     c.ethApi.auth.From,
    Signer:   c.ethApi.auth.Signer,
    GasLimit: 2381623,
  }, address)
  if err != nil {
    log.Println("Could not accept ride request from car: %v", err)
  }
  fmt.Println("Transaction initiated")
  receipt, err := bind.WaitMined(context.Background(), c.ethApi.conn, transaction)
  if err != nil {
    log.Fatalf("Wait for mining error %s %v: ", err)
  } else if receipt.Status == types.ReceiptStatusFailed {
    log.Printf("Accept Request failed \n")
    c.acceptanceState = Fail
    fmt.Println("Transaction Failed")
  } else {
    log.Printf("Accept Request success \n")
    c.ethApi.chosenRider = address
    c.acceptanceState = Success
    fmt.Println("Transaction Success")
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

