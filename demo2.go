package main

/*
import (
	"github.com/gorilla/websocket"
	"net/http"
	"log"
	"demo2/go-packages/simulator"
	"strconv"
)

var clients = make(map[*websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)           // broadcast channel

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Define our message object
type Message struct {
	Type          string `json:"type"`
	Id            string `json:"id"`
	X             string `json:"x"`
	Y             string `json:"y"`
	Orientation   string `json:"orientation"`
	State         string `json:"state"`
}

func main() {
	// Create a simple file server
	fs := http.FileServer(http.Dir("public"))
	http.Handle("/", fs)

	// Configure websocket route
	http.HandleFunc("/ws", handleConnections)

	// Upgrade HTTP to WS and store client handle
	go handleMessages()

	go func() {
		// Start the server on localhost portNo and log any errors
		log.Printf("http server started on %s \n", ":8000")
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}()

	world := simulator.GetWorldFromFile("maps/4by4.map")
	world.SetFrameRate(50)
	world.SetUnitsPerMove(1)
	world.AddCars(4)
	go world.Loop()
	receiveChannel := world.GetBroadcastChannel()
	for {
		carStates := <-receiveChannel
		for _, v := range carStates {
			//fmt.Println(v.Id, " ", v.Coordinates, " ", v.Orientation)
			broadcast <- Message{Type:"Car", Id:strconv.Itoa(int(v.Id)), X:strconv.Itoa(int(v.Coordinates.X)),
									Y:strconv.Itoa(int(v.Coordinates.Y)), Orientation:strconv.Itoa(int(v.Orientation))}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg Message
		// Read in a new message as JSON and map it to a Message object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
	}
}

func handleMessages() {
	for {
		// Grab the next message from the broadcast channel
		msg := <-broadcast
		// Send it out to every client that is currently connected
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
//*/

// Section for testing
//*
import (
  "demo2/go-packages/sim2"
  "log"
  "fmt"
  "os"
  "bufio"
)

// TODO: remove commented-out test prints and make proper test files

func main() {
  fmt.Println("Starting demo2 simulation")

  // Instantiate world
  w := sim2.GetWorldFromFile("maps/4by4.map")
  w.Fps = float64(50)

  keys, err := os.Open("keys.txt")
  if err != nil {
    log.Fatal(err)
  }
  defer keys.Close()

  scanner := bufio.NewScanner(keys)

  scanner.Scan()
  scanner.Scan()
  //address of the deployed ferris contract
  existingMrmAddress := scanner.Text()

  // Instantiate cars
  numCars := uint(1)
  cars := make([]*sim2.Car, numCars)
  for i := uint(0); i < numCars; i++ {
    // Request to register new car from World
    syncChan, updateChan, ok := w.RegisterCar(i)
    if !ok {
      log.Printf("error: failed to register car")
    }

    scanner.Scan()
    scanner.Scan()
    carPrivateKey := scanner.Text()
    // Construct car
    cars[i] = sim2.NewCar(i, w, syncChan, updateChan, existingMrmAddress, carPrivateKey)
  }

  // Instantiate JSON web output
  webChan, ok := w.RegisterWeb()
  if !ok {
    log.Printf("error: failed to register web output")
  }
  web := sim2.NewWebSrv(webChan)

  // Begin World operation
  go w.LoopWorld()

  // Begin Car operation
  for i := uint(0); i < numCars; i++ {
    go cars[i].CarLoop()
  }

  // Begin JSON web output operation
  go web.LoopWebSrv()

  select{}  // Do work in the coroutines, main has nothing left to do
}
//*/
