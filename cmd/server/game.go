package main

import (
	"fmt"
	"sync"
)

type board [9]string

type game struct {
	board    board
	turn     string
	finished bool
	mu       sync.Mutex
}

func newGame() *game {
	return &game{turn: "X"}
}

const (
	ongoing = "ongoing"
	draw    = "draw"
)

var winningLines = [8][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8},
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8},
	{0, 4, 8}, {2, 4, 6},
}

// play validates and records one move. The mutex protects the shared game
// because both client handlers can submit moves concurrently.
// play validates and records one move, returning the resulting game status.
func (g *game) play(symbol string, cell int) (string, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	if err := g.validateMove(symbol, cell); err != nil {
		return "", err
	}

	g.board[cell-1] = symbol
	status := checkOutcome(g.board)
	if status != ongoing {
		g.finished = true
		return status, nil
	}
	if g.turn == "X" {
		g.turn = "O"
	} else {
		g.turn = "X"
	}
	return ongoing, nil
}

func checkOutcome(board board) string {
	for _, line := range winningLines {
		symbol := board[line[0]] // get the symbol in the first cell of the line

		// if symbol non empty, and the other two cells in the line match, we have a winner
		if symbol != "" && symbol == board[line[1]] && symbol == board[line[2]] {
			return symbol
		}
	}

	// check if any cell is empty; if so, the game is still ongoing
	for _, cell := range board {
		if cell == "" {
			return ongoing
		}
	}
	return draw
}

func (g *game) snapshot() (board, string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.board, g.turn
}

// validateMove checks every rule that must be true before a move is applied.
// The caller must hold g.mu.
func (g *game) validateMove(symbol string, cell int) error {
	if g.finished {
		return fmt.Errorf("game is over")
	}
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
