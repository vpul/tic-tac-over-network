package client

import (
	"fmt"
	"strconv"
)

func (c *Client) renderBoard() {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.outputMu.Lock()
	defer c.outputMu.Unlock()
	c.renderBoardLocked()
}

func (c *Client) renderBoardLocked() {
	fmt.Fprintln(c.output)
	for r := 0; r < 3; r++ {
		i := r * 3
		cells := [3]string{c.board[i], c.board[i+1], c.board[i+2]}
		for j := range cells {
			if cells[j] == "" {
				cells[j] = strconv.Itoa(i + j + 1)
			}
		}
		fmt.Fprintf(c.output, " %s | %s | %s \n", cells[0], cells[1], cells[2])
		if r < 2 {
			fmt.Fprintln(c.output, "---+---+---")
		}
	}
	fmt.Fprintln(c.output, "\nplay a cell (1-9) or q to quit:")
}
