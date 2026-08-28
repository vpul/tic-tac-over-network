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
			if !c.isPaired() {
				c.println("not paired yet — wait for an opponent")
				continue
			}
			if err := encoder.Encode(protocol.Request{Type: "move", Cell: cell}); err != nil {
				return fmt.Errorf("send move: %w", err)
			}
		}
	}
	return c.input.Err()
}

func (c *Client) isPaired() bool {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	return c.symbol != ""
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
