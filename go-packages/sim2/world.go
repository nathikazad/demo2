package sim2

import (
  "time"
  "strconv"
  //"fmt"
)

// world - Describes world state: position and velocity of all cars in simulation

// CarInfo - struct to contain position and velocity information for a simulated car.
type CarInfo struct {
  ID uint
  Pos Coords
  Vel Coords  // with respect to current position, offset for a single frame
  Dir Coords  // unit vector with respect to current position
}

// World - struct to contain all relevat world information in simulation.
type World struct {
  Graph *Digraph
  CarStates map[uint]CarInfo
  Fps float64

  syncChans map[uint]chan bool  // Index by actor ID for channel to/from that actor
  recvChans map[uint]chan CarInfo  // Index by actor ID for channel to/from that actor
  webChan chan Message
}

// NewWorld - Constructor for valid World object.
func NewWorld() *World {
  w := new(World)
  w.Graph = NewDigraph()
  w.CarStates = make(map[uint]CarInfo)

  w.Fps = float64(1)
  w.syncChans = make(map[uint]chan bool)
  w.recvChans = make(map[uint]chan CarInfo)
  // NOTE webChan is nil until registered
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
  w.CarStates[ID] = CarInfo{ ID:ID }  // TODO: randomize/control car location on startup
  w.syncChans[ID] = make(chan bool, 1)  // Buffer up to one output
  w.recvChans[ID] = make(chan CarInfo, 1)  // Buffer up to one input
  return w.syncChans[ID], w.recvChans[ID], true
}

// TODO: an UnregisterCar func if necessary

// RegisterWeb - If not already registered, allocate a channel for web output and true OK.
func (w *World) RegisterWeb() (chan Message, bool) {
  // Check for invalid world or previous allocation
  if w == nil || w.webChan != nil{
    return nil, false
  }

  // Allocate new channel for registered web output
  w.webChan = make(chan Message)
  return w.webChan, true
}

// TODO: an UnregisterWeb func if necessary

// LoopWorld - Begin the world simulation execution loop
func (w *World) LoopWorld() {
  counter := uint64(0)
  for {
    //fmt.Println("Iteration", counter)
    timer := time.NewTimer(time.Duration(1000/w.Fps) * time.Millisecond)

    // Send out sync flag = true for each registered car
    for ID := range w.CarStates {
      w.syncChans[ID] <- true
    }

    // Car coroutines should now process current world state
    for _, car := range w.CarStates {
      w.webChan <- Message{
        Type:"Car",
        ID:strconv.Itoa(int(car.ID)),
        X:strconv.Itoa(int(car.Pos.X)),
        Y:strconv.Itoa(int(car.Pos.Y)),
        Orientation:strconv.Itoa(int(Coords{0,0}.Angle(car.Dir))), 
      }
    }

    // Set up routines to catch incoming car data
    recvDone := make(chan bool, len(w.recvChans))
    for ID, carChan := range w.recvChans {
      go func(idx uint, dataChan chan CarInfo, dstMap map[uint]CarInfo) {
        data := <- dataChan
        //fmt.Println("World got new data on index", idx, ":", data)
        // TODO: any validation of data here
        // TODO: also consider a clone method for CarInfo for future-proofing copy assign
        dstMap[idx] = data
        recvDone <- true
        //fmt.Println("World recvDone sent on", idx)
      }(ID, carChan, w.CarStates)
    }

    // Wait for all channels to report
    for i := 0; i < len(w.recvChans); i++ {
      <-recvDone
      //fmt.Println("World recvDone recv on", i)
    }
    //fmt.Println("World done getting all recvDone")

    counter++

    // Wait for frame update
    <-timer.C

    // World loop iterates
  }
}
