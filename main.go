package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"sync"
)

// ClientManager is a data structure on the server that helps keep track of
// connected clients. The broadcast channel is used to broadcast received
// messages to all clients. Register and unregister are to correctly
// log and handle the connection or disconnection of clients.
type ClientManager struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

// Node
type Node struct {
	serverPort string
	id         int
	name       string
	color      string
	manager    ClientManager
}

// Client holds the connection information for each client. The data channel is
// used to send and receive messages as raw bytes.
type Client struct {
	socket net.Conn
	name   string
	color  string
	data   chan []byte
}

// Constants
var PROTOCOL = "tcp"
var SERVER_BASE_PORT = 10000

func main() {
	flagID := flag.Int("id", 1, "This node's id number")
	flag.Parse()
	serverPort := strconv.Itoa(SERVER_BASE_PORT + *flagID)

	fmt.Print("Enter name: ")
	reader := bufio.NewReader(os.Stdin)
	name, _ := reader.ReadString('\n')

	fmt.Print("Choose color (red/blue/green): ")
	color, _ := reader.ReadString('\n')

	node := Node{
		serverPort: serverPort,
		id:         *flagID,
		name:       name[:len(name)-2],
		color:      color,
	}
	go node.startServer()
	go node.startClient()

	fmt.Println("Launched node. Waiting for connections...")

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}
