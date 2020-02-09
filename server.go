package main

import (
	"fmt"
	"net"
	"strings"

	"github.com/fatih/color"
)

// Starts the ClientManager. It defines how to handle inputs on each of
// the ClientManager's channels.
func (manager *ClientManager) start(node *Node) {
	for {
		select {
		// Receiving a request to register a new client
		case client := <-manager.register:
			manager.clients[client] = true
			fmt.Printf("[+] New connection from %s. Total connected: %d\n", client.socket.RemoteAddr(),
				len(manager.clients))
			client.socket.Write([]byte("/resetname " + node.name + "\n"))
			client.socket.Write([]byte("/setcolor " + node.color + "\n"))
			// Receiving a request to unregister a client
		case client := <-manager.unregister:
			if _, ok := manager.clients[client]; ok {
				close(client.data)
				delete(manager.clients, client)
				fmt.Printf("[-] %s disconnected. Total connected: %d\n", client.socket.RemoteAddr(), len(manager.clients))
			}
			// Receiving message on the broadcast channel => send it to all clients
		case message := <-manager.broadcast:
			for client := range manager.clients {
				// Write the message to each client's data channel, so it will be
				// sent out by server send()
				select {
				case client.data <- message:
				default:
					close(client.data)
					delete(manager.clients, client)
				}
			}
		}
	}
}

// The server receive function, launched from the startServer function.
// Listens for incoming messages from clients, and broadcasts received
// messages to all other clients.
func (manager *ClientManager) receive(client *Client, node *Node) {
	for {
		message := make([]byte, 4096)
		length, err := client.socket.Read(message)

		if err != nil {
			// If received an error on Read, unregister the client
			manager.unregister <- client
			client.socket.Close()
			break
		}

		if length > 0 {
			// Handle commands
			messageStr := string(message[:length-1])
			messageSplit := strings.SplitN(messageStr, " ", 3)
			if messageSplit[0] == "/resetname" {
				client.name = messageSplit[1]
			} else if messageSplit[0] == "/msg" {
				if messageSplit[1] == node.name {
					output := client.name + " (whisper): " + messageSplit[2]
					color.Yellow(output)
				}
			} else if messageSplit[0] == "/setcolor" {
				client.color = messageSplit[1]
			} else {
				// message[length] is Carriage Return since this is on Windows
				output := client.name + ": " + messageStr
				c := getColor(client.color[:len(client.color)-2])
				c.Println(output)
			}
		}
	}
}

// The server send function launched as a goroutine. It simply watches for
// data/messages written to this client's data channel, then writes the message
// to the client's socket. This send out the message to the client.
func (manager *ClientManager) send(client *Client) {
	defer client.socket.Close()
	for {
		select {
		case message, ok := <-client.data:
			if !ok {
				return
			}
			client.socket.Write(message)
		}
	}
}

// The server start function. Launches the server on the specified port, and
// sets up the ClientManager.
func (node *Node) startServer() {
	listener, error := net.Listen(PROTOCOL, ":"+node.serverPort)
	if error != nil {
		fmt.Println(error)
	}
	node.manager = ClientManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
	go node.manager.start(node)

	node.connectAllNodes()

	for {
		// Blocks until a new connection arrives
		connection, _ := listener.Accept()
		if error != nil {
			fmt.Println(error)
		}
		client := &Client{socket: connection, data: make(chan []byte), name: "guest_" + connection.RemoteAddr().String()}
		node.manager.register <- client
		go node.manager.receive(client, node)
		go node.manager.send(client)
	}
}
