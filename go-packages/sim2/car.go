package sim2

import (
  //"fmt"
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
