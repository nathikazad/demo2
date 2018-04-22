package simulator

type Car struct {
	id uint
	path []uint
	currentState CarState
}

type CommandAction int
const (
	Move CommandAction = 0
	Stop CommandAction = 1
	Turn CommandAction = 2
	Park CommandAction = 3
)

type Action int
const (
	Moving Action = 0
	Stopped Action = 1
	Turning Action = 2
	Parked Action = 3
)

type CarCommand struct {
	id uint
	commandAction CommandAction
	arg0 uint
}

type CarState struct {
	id uint
	coordinates Coordinates
	orientation uint8
	edgeId uint
	action Action
}

func (c *Car) startLoop(broadcast chan map[*Car]CarState,
					commandReceiver chan CarCommand) {
	for {
		carStates := <-broadcast
		c.currentState = carStates[c]
		commandReceiver <- CarCommand{c.id, Stop, 0}
	}
}