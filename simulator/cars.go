package simulator

type Car struct {
	id uint
	path []uint
	currentState CarState
	StateReceiver chan map[*Car]CarState
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
	Turning PresentAction = 2
	Parked PresentAction = 3
)

type CarCommand struct {
	id uint
	commandAction CommandAction
	arg0 uint
}

type CarState struct {
	Id uint
	Coordinates Coordinates
	Orientation uint
	edgeId uint
	presentAction PresentAction
}

func (c *Car) startLoop(commandSender chan<- CarCommand) {
	c.StateReceiver = make(chan map[*Car]CarState)
	go func() {
		for {
			carStates := <- c.StateReceiver
			c.currentState = carStates[c]
			commandSender <- c.produceNextCommand()
		}
	}()
}

func (c *Car) produceNextCommand() CarCommand {
	return CarCommand{c.id, Move, 0}
}