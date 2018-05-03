package sim2

import (
  "github.com/gorilla/websocket"
  "net/http"
	"log"
  //"fmt"
)

// Message - struct to contain all data relevant to rendering a Car on the frontend.
type Message struct {
	Type          string `json:"type"`
	ID            string `json:"id"`
	X             string `json:"x"`
	Y             string `json:"y"`
	Orientation   string `json:"orientation"`
	State         string `json:"state"`
}

// WebSrv - container for web server variables for the simulator.
type WebSrv struct {
  // TODO: other fields here as necessary
  webChan chan Message  // Incoming car information from simulator
}

// TODO: find a way to move these into the WebSrv struct without violating handleConnection
var clients = make(map[*websocket.Conn]bool) // connected clients
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewWebSrv - Constructor for a valid WebSrv object.
func NewWebSrv(web chan Message) *WebSrv {
  s := new(WebSrv)
  s.webChan = web
  return s
}

// LoopWebSrv - Begin the web server execution loop.
func (s *WebSrv) LoopWebSrv() {
  // Create a simple file server
  fs := http.FileServer(http.Dir("public"))
  http.Handle("/", fs)

  // Configure websocket route
  http.HandleFunc("/ws", handleConnections)

  // Start the server on localhost portNo and log any errors
  go func() {
		log.Printf("http server started on %s \n", ":8000")
		err := http.ListenAndServe(":8000", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
  }()

  // Handle any car info updates from World
  for {
    // Grab the next message from the broadcast channel
    msg := <-s.webChan
    //fmt.Println("Got msg:",msg)
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
