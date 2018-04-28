package simulator

import (
	"math/rand"
	"time"
)

type Car struct {
	id uint
	path []uint
	currentState CarState
	StateReceiver chan map[*Car]CarState
	world *World
}

type CommandAction int
const (
	Move CommandAction = 0
	Stop CommandAction = 1
	Turn CommandAction = 2
	Park CommandAction = 3
)

type PresentAction int
const (
	Moving PresentAction = 0
	Stopped PresentAction = 1
	StoppedAtTurn PresentAction = 2
	Turning PresentAction = 3
	Parked PresentAction = 4
)

type CarCommand struct {
	id uint
	commandAction CommandAction
	arg0 uint
}

type CarState struct {
	Id uint
	Coordinates Coordinates
	Orientation float64
	edgeId uint
	presentAction PresentAction
}

func (c *Car) startLoop(commandSender chan<- CarCommand) {
	c.StateReceiver = make(chan map[*Car]CarState)
	go func() {
		for {
			carStates := <- c.StateReceiver
			c.currentState = carStates[c]
			commandSender <- c.produceNextCommand(carStates)
		}
	}()
}

func (c *Car) produceNextCommand(map[*Car]CarState) CarCommand {
	switch c.currentState.presentAction {
	case Moving:
		// Check for collision
		return CarCommand{c.id, Move, 0}
	case Stopped:
		return CarCommand{c.id, Move, 0}
	case StoppedAtTurn:
		_, endNodeId := c.world.GetStartAndEndNodeOfEdge(c.currentState.edgeId)
		adjacentNodeIds := c.world.GetAdjacentNodeIds(endNodeId)
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)
		choosenNodeId := adjacentNodeIds[uint(r1.Int() % len(adjacentNodeIds))]
		return CarCommand{id:c.id, commandAction:Turn, arg0:choosenNodeId}
	default :
		return CarCommand{id:c.id, commandAction:Turn, arg0:0}
	}

}