package client

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"tic-tac-over-network/internal/protocol"
)

func (c *Client) processInput() error {
	encoder := json.NewEncoder(c.conn)
	for c.input.Scan() {
		line := strings.TrimSpace(c.input.Text())
		switch {
		case line == "":
			continue
		case line == "q" || line == "quit":
			return nil
		default:
			cell, err := strconv.Atoi(line)
			if err != nil || cell < 1 || cell > 9 {
				c.println("? enter a cell 1-9, or q to quit")
				continue
			}
			if !c.previewMove(cell) {
				continue
			}
			if err := encoder.Encode(protocol.Request{Type: "move", Cell: cell}); err != nil {
				return fmt.Errorf("send move: %w", err)
			}
		}
	}
	return c.input.Err()
}

func (c *Client) previewMove(cell int) bool {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	if c.symbol == "" {
		c.println("not paired yet — wait for an opponent")
		return false
	}
	if c.turn != c.symbol {
		c.printf("not your turn — waiting for %s\n", c.turn)
		return false
	}
	if c.board[cell-1] != strconv.Itoa(cell) {
		c.printf("(local preview) cell %d already marked %s — not sent\n", cell, c.board[cell-1])
		return false
	}
	c.board[cell-1] = c.symbol
	c.renderBoardLocked()
	return true
}

func (c *Client) println(value string) {
	c.outputMu.Lock()
	defer c.outputMu.Unlock()
	fmt.Fprintln(c.output, value)
}

func (c *Client) printf(format string, args ...any) {
	c.outputMu.Lock()
	defer c.outputMu.Unlock()
	fmt.Fprintf(c.output, format, args...)
}
