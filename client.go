package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

// Attempt to connect to all launched nodes based on id number.
func (node *Node) connectAllNodes() {
	possibleNodes := node.getPossibleNodes()
	for _, port := range possibleNodes {
		connection, error := net.Dial(PROTOCOL, "localhost:"+port)
		if error != nil {
			fmt.Errorf("%s: No node at port %s", error.Error(), port)
			continue
		}
		client := &Client{socket: connection, data: make(chan []byte), name: "guest_" + node.serverPort}
		node.manager.register <- client
		go node.manager.receive(client, node)
		go node.manager.send(client)
	}
}

// The client start function. Sets up the socket with the specified port, and
// gets ready to take input from stdin.
func (node *Node) startClient() {
	for {
		reader := bufio.NewReader(os.Stdin)
		message, _ := reader.ReadString('\n')

		// Handle commands
		messageSplit := strings.Split(message, " ")
		if messageSplit[0] == "/resetname" {
			node.name = messageSplit[1][:len(messageSplit[1])-2]
		}

		// Write to socket
		for client, _ := range node.manager.clients {
			client.socket.Write([]byte(strings.TrimRight(message, "\n")))
		}
	}
}
