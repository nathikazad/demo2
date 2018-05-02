package sim2

import (
  "time"
  "fmt"
)

// world - Describes world state: position and velocity of all cars in simulation

// CarInfo - struct to contain position and velocity information for a simulated car.
type CarInfo struct {
  Pos Coords
  Vel Coords  // with respect to current position, offset for a single frame
}

// World - struct to contain all relevat world information in simulation.
type World struct {
  Graph *Digraph
  CarStates map[uint]CarInfo
  LastTime time.Time

  fps float64
  syncChans map[uint]chan bool  // Index by actor ID for channel to/from that actor
  recvChans map[uint]chan CarInfo  // Index by actor ID for channel to/from that actor
}

// NewWorld - Constructor for valid World object.
func NewWorld() *World {
  w := new(World)
  w.Graph = NewDigraph()
  w.CarStates = make(map[uint]CarInfo)

  w.fps = float64(1)
  w.syncChans = make(map[uint]chan bool)
  w.recvChans = make(map[uint]chan CarInfo)
  return w
}

// GetWorldFromFile - Populate underlying diggraph for roads on world.
func GetWorldFromFile(fname string) (w *World) {
  w = NewWorld()
  w.Graph = GetDigraphFromFile(fname)
  return w
}

// RegisterCar - If the car ID has not been taken, allocate new channels for the car ID and true OK.
//   If the car ID is taken or World is unallocated, return nil channels and false OK value.
func (w *World) RegisterCar(ID uint) (chan bool, chan CarInfo, bool) {
  // Check for invalid World
  if w == nil || w.syncChans == nil || w.recvChans == nil {
    return nil, nil, false
  }

  // Check for previous allocation
  if _, ok := w.CarStates[ID]; ok {
    return nil, nil, false
  }

  // Allocate new channels for registered car
  w.CarStates[ID] = CarInfo{}
  w.syncChans[ID] = make(chan bool, 1)  // Buffer up to one output
  w.recvChans[ID] = make(chan CarInfo, 1)  // Buffer up to one input
  return w.syncChans[ID], w.recvChans[ID], true
}

// TODO: an UnregisterCar func if necessary

// LoopWorld - Begin the world simulation execution loop
func (w *World) LoopWorld() {
  w.LastTime = time.Now()
  counter := uint64(0)
  for {
    fmt.Println("Iteration", counter)
    w.LastTime = time.Now()
    timer := time.NewTimer(time.Duration(1000/w.fps) * time.Millisecond)

    // Send out sync flag = true for each registered car
    for ID := range w.CarStates {
      w.syncChans[ID] <- true
    }

    // Car coroutines should now process current world state
    // TODO: JSON output should be written now, need a different sync chan to inform it

    // Set up routines to catch incoming car data
    recvDone := make(chan bool, len(w.recvChans))
    for ID, carChan := range w.recvChans {
      go func(idx uint, dataChan chan CarInfo, dstMap map[uint]CarInfo) {
        data := <- dataChan
        fmt.Println("World got new data on index", idx, ":", data)
        dst := dstMap[idx]
        dst.Pos = data.Pos
        dst.Vel = data.Vel
        recvDone <- true
        fmt.Println("World recvDone sent on", idx)
      }(ID, carChan, w.CarStates)
    }

    // Wait for all channels to report
    for i := 0; i < len(w.recvChans); i++ {
      <-recvDone
      fmt.Println("World recvDone recv on", i)
    }
    fmt.Println("World done getting all recvDone")

    counter++

    // Wait for frame update
    <-timer.C

    // World loop iterates
  }
}
