package server

import (
	"fmt"
	"sync"

	"tic-tac-over-network/internal/game"
	"tic-tac-over-network/internal/protocol"
)

// session owns one matched game's workflow and participants.
type session struct {
	mu      sync.Mutex
	game    *game.Game
	clients [2]*client
	closed  bool
}

func newSession(first, second *client) *session {
	return &session{
		game:    game.New(),
		clients: [2]*client{first, second},
	}
}

// start assigns game roles and makes the session available to both clients.
func (s *session) start() {
	board, turn := s.game.Snapshot()
	for i, symbol := range []string{"X", "O"} {
		client := s.clients[i]
		client.session = s
		client.send(protocol.Response{
			Type:   "game_start",
			Symbol: symbol,
			Board:  board,
			Turn:   turn,
		})
		fmt.Printf("%s assigned %s\n", client.remote(), symbol)
	}
	for _, client := range s.clients {
		close(client.sessionReady)
	}
}

func (s *session) symbolFor(client *client) (string, bool) {
	switch client {
	case s.clients[0]:
		return "X", true
	case s.clients[1]:
		return "O", true
	default:
		return "", false
	}
}

func (s *session) broadcast(message protocol.Response) {
	for _, client := range s.clients {
		client.send(message)
	}
}

// handle serializes validation, state mutation, and broadcasts for this game.
func (s *session) handle(client *client, p payload) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.closed {
		client.send(protocol.Response{Type: "error", Reason: "game is over"})
		return
	}
	symbol, ok := s.symbolFor(client)
	if !ok {
		client.send(protocol.Response{Type: "error", Reason: "client is not part of this session"})
		return
	}

	switch p.message.Type {
	case "move":
		status, err := s.game.Play(symbol, p.message.Cell)
		if err != nil {
			client.send(protocol.Response{Type: "error", Reason: err.Error()})
			return
		}

		board, turn := s.game.Snapshot()
		if status != game.Ongoing {
			s.broadcast(protocol.Response{Type: "game_over", Board: board, Turn: turn, Result: status})
			fmt.Printf("game over: %s\n", status)
			return
		}

		s.broadcast(protocol.Response{Type: "state", Board: board, Turn: turn})
		fmt.Printf("move from %s: %s played cell %d\n", client.remote(), symbol, p.message.Cell)
	default:
		client.send(protocol.Response{Type: "error", Reason: "unknown message type"})
		fmt.Printf("ignored message from %s: unknown type %q\n", client.remote(), p.message.Type)
	}
}

func (s *session) disconnect(client *client) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	s.closed = true
	for _, opponent := range s.clients {
		if opponent != client {
			opponent.send(protocol.Response{Type: "opponent_left"})
		}
	}
}
