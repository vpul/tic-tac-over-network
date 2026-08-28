package client

import (
	"bufio"
	"io"
	"net"
	"sync"
)

// Client coordinates the network reader and terminal input for one player.
type Client struct {
	conn   net.Conn
	input  *bufio.Scanner
	output io.Writer

	stateMu  sync.Mutex
	outputMu sync.Mutex
	symbol   string
	turn     string
	board    [9]string
}

func New(conn net.Conn, input io.Reader, output io.Writer) *Client {
	return &Client{
		conn:   conn,
		input:  bufio.NewScanner(input),
		output: output,
		turn:   "X",
		board:  [9]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"},
	}
}

// Run starts the server reader and processes terminal commands until exit.
func (c *Client) Run() error {
	go c.readServer()
	c.renderBoard()
	return c.processInput()
}
