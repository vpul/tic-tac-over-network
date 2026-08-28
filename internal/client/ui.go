package client

import (
	"fmt"
	"strconv"
)

func (c *Client) renderBoard() {
	c.stateMu.Lock()
	symbol := c.symbol
	turn := c.turn
	board := c.board
	c.stateMu.Unlock()

	c.outputMu.Lock()
	defer c.outputMu.Unlock()

	if symbol != "" {
		fmt.Fprintf(c.output, "your symbol: %s\n", symbol)
	}

	if turn != "" {
		fmt.Fprintf(c.output, "current turn: %s\n", turn)
	}

	fmt.Fprintln(c.output)

	for row := 0; row < 3; row++ {
		rowStart := row * 3
		// Copy the current row's values from the board.
		cells := []string{board[rowStart], board[rowStart+1], board[rowStart+2]}

		// Replace empty cells with selectable cell numbers.
		for column := range cells {
			if cells[column] == "" {
				cells[column] = strconv.Itoa(rowStart + column + 1)
			}
		}
		fmt.Fprintf(c.output, " %s | %s | %s \n", cells[0], cells[1], cells[2])
		if row < 2 {
			fmt.Fprintln(c.output, "---+---+---")
		}
	}
	fmt.Fprintln(c.output, "\nplay a cell (1-9) or q to quit:")
}
