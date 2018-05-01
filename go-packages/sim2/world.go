package sim2

import (
  "time"
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

  dtS float64
  fps float64
  lastTime time.Time
  syncChans map[uint]chan bool  // Index by actor ID for channel to/from that actor
  recvChans map[uint]chan CarInfo  // Index by actor ID for channel to/from that actor
}

// NewWorld - Constructor for valid World object.
func NewWorld() *World {
  w := new(World)
  w.Graph = NewDigraph()
  w.CarStates = make(map[uint]CarInfo)

  w.dtS = float64(0)
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

// TODO: UnregisterCar if necessary

// GetDtS - Fetch time duration since last update of the simulation, in [s].
func (w *World) GetDtS() float64 {
  return w.dtS
}

// LoopWorld - Begin the world simulation execution loop
func (w *World) LoopWorld() {
  w.lastTime = time.Now()
  for {
    // Calculate dt
    w.dtS = time.Since(w.lastTime).Seconds()
    w.lastTime = time.Now()

    // Send out sync flag = true for each registered car
    for ID := range w.CarStates {
      w.syncChans[ID] <- true
    }

    // Car coroutines should now process current world state

    // Wait for all car data or frame refresh
    // TODO continue here
  }
}

// TODO JSON output externally, need a new channel exposed to sync that operation as well
