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
  "fmt"
  "time"
)

func main() {
  w := sim2.GetWorldFromFile("maps/4by4.map")
  fmt.Println(w)

  ID := uint(0)
  syncChan, updateChan, ok0 := w.RegisterCar(ID)
  fmt.Println("Register '0' #1 successful?", ok0)
  _, _, ok1 := w.RegisterCar(ID)
  fmt.Println("Register '0' #2 successful?", ok1)

  go func(world *sim2.World, idx uint, sync chan bool, update chan sim2.CarInfo) {
    for {
      <-sync
      fmt.Println("Car", ID, ": got sync")
      fmt.Println("Car", ID, ": DT - ", time.Since(w.LastTime))
      update <- sim2.CarInfo{sim2.Coords{1,2}, sim2.Coords{3,4}}
      fmt.Println("Car", ID, ": sent empty car info")
    }
  }(w, ID, syncChan, updateChan)

  go w.LoopWorld()

  for {
    <-chan bool (nil)
  }  // Do work in the coroutines
}
//*/
