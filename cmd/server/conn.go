package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
)

// handleConn owns one connection for its whole lifetime: parks it in the
// lobby, waits for the matchmaker to assign a symbol, then keeps the conn
// open until the client disconnects.
func handleConn(conn net.Conn) {
	defer conn.Close()
	remote := conn.RemoteAddr()
	fmt.Printf("client connected: %s\n", remote)
	defer fmt.Printf("client disconnected: %s\n", remote)

	connectedClient := &client{
		conn:         conn,
		symbol:       make(chan string, 1),
		messages:     make(chan message),
		disconnected: make(chan struct{}),
	}

	// Decode messages from the moment we connect, so a lobby disconnect
	// is visible to the matchmaker via connectedClient.disconnected.
	go func() {
		defer close(connectedClient.disconnected)
		dec := json.NewDecoder(conn)
		for {
			var m message
			if err := dec.Decode(&m); err != nil {
				if !errors.Is(err, io.EOF) {
					fmt.Printf("read error from %s: %v\n", remote, err)
				}
				return
			}
			connectedClient.messages <- m
		}
	}()

	waitingClients <- connectedClient // park in the lobby
	assignedSymbol := <-connectedClient.symbol
	fmt.Printf("%s assigned %s\n", remote, assignedSymbol)

	for {
		select {
		case m := <-connectedClient.messages:
			switch m.Type {
			case "move":
				fmt.Printf("move from %s: cell %d\n", remote, m.Cell)
			default:
				fmt.Fprintf(conn, `{"type":"error","reason":"unknown message type"}`+"\n")
				fmt.Printf("ignored message from %s: unknown type %q\n", remote, m.Type)
			}
		case <-connectedClient.disconnected:
			return
		}
	}
}
