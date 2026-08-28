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
	board    [9]string
}

func New(conn net.Conn, input io.Reader, output io.Writer) *Client {
	return &Client{
		conn:   conn,
		input:  bufio.NewScanner(input),
		output: output,
	}
}

// Run starts the server reader and processes terminal commands until exit.
func (c *Client) Run() error {
	go c.readServer()
	c.renderBoard()
	return c.processInput()
}
