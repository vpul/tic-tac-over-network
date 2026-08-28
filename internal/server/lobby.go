package server

import (
	"fmt"

	"tic-tac-over-network/internal/protocol"
)

// runLobby pairs clients first-come-first-served.
func (s *Server) runLobby() {
	for {
		waitingClient := <-s.waitingClients
		waitingClient.send(protocol.Response{Type: "waiting"})

		opponent, ok := s.waitForOpponent(waitingClient)
		if !ok {
			continue
		}

		matchedSession := newSession(waitingClient, opponent)
		matchedSession.start()
		fmt.Printf("paired %s (X) vs %s (O)\n", waitingClient.remote(), opponent.remote())
	}
}

func (s *Server) waitForOpponent(waitingClient *client) (*client, bool) {
	select {
	case opponent := <-s.waitingClients:
		return opponent, true
	case <-waitingClient.disconnected:
		fmt.Printf("%s left the lobby\n", waitingClient.remote())
		return nil, false
	}
}
