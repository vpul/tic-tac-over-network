package client

import "fmt"

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
		fmt.Fprintf(c.output, " %s | %s | %s \n", c.board[i], c.board[i+1], c.board[i+2])
		if r < 2 {
			fmt.Fprintln(c.output, "---+---+---")
		}
	}
	fmt.Fprintln(c.output, "\nplay a cell (1-9) or q to quit:")
}
