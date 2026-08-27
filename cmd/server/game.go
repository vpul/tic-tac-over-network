package main

import (
	"fmt"
	"sync"
)

type game struct {
	board [9]string
	turn  string
	mu    sync.Mutex
}

func newGame() *game {
	return &game{turn: "X"}
}

// play validates and records one move. The mutex protects the shared game
// because both client handlers can submit moves concurrently.
func (g *game) play(symbol string, cell int) error {
	g.mu.Lock()
	defer g.mu.Unlock()

	if err := g.validateMove(symbol, cell); err != nil {
		return err
	}

	g.board[cell-1] = symbol
	if g.turn == "X" {
		g.turn = "O"
	} else {
		g.turn = "X"
	}
	return nil
}

func (g *game) snapshot() (board [9]string, turn string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.board, g.turn
}

// validateMove checks every rule that must be true before a move is applied.
// The caller must hold g.mu.
func (g *game) validateMove(symbol string, cell int) error {
	if cell < 1 || cell > len(g.board) {
		return fmt.Errorf("cell must be between 1 and %d", len(g.board))
	}
	if symbol != g.turn {
		return fmt.Errorf("it is %s's turn", g.turn)
	}
	if g.board[cell-1] != "" {
		return fmt.Errorf("cell %d is already occupied", cell)
	}
	return nil
}
