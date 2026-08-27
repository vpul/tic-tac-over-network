package main

import (
	"fmt"
	"net"
)

// client represents one connected client waiting for an opponent.
type client struct {
	conn         net.Conn
	symbol       chan string   // matchmaker sends "X" or "O" once paired
	messages     chan message  // reader sends client messages here
	disconnected chan struct{} // closed when the connection dies
}

// waitingClients is the lobby queue. Its capacity of 1 lets the first client
// park here until a second client arrives and the matchmaker pairs them.
var waitingClients = make(chan *client, 1)

// matchmaker pairs clients first-come-first-served. The first client plays X;
// the second client plays O.
func matchmaker() {
	for {
		waitingClient := <-waitingClients
		fmt.Fprintln(waitingClient.conn, `{"type":"waiting"}`)

		opponent, ok := waitForOpponent(waitingClient)
		if !ok {
			continue // waiting client left; find a new waiting client
		}

		fmt.Fprintln(waitingClient.conn, `{"type":"game_start","symbol":"X"}`)
		fmt.Fprintln(opponent.conn, `{"type":"game_start","symbol":"O"}`)
		waitingClient.symbol <- "X"
		opponent.symbol <- "O"
		fmt.Printf("paired %s (X) vs %s (O)\n", waitingClient.conn.RemoteAddr(), opponent.conn.RemoteAddr())
	}
}

// waitForOpponent blocks until another client arrives. If the waiting client
// disconnects first, it returns ok=false instead of pairing a dead connection.
func waitForOpponent(waitingClient *client) (opponent *client, ok bool) {
	select {
	case opponent := <-waitingClients:
		return opponent, true
	case <-waitingClient.disconnected:
		fmt.Printf("%s left the lobby\n", waitingClient.conn.RemoteAddr())
		return nil, false
	}
}
