package client

import (
	"bufio"
	"io"
	"net"
	"sync"

	"tic-tac-over-network/internal/game"
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
	board    game.Board
}

func New(conn net.Conn, input io.Reader, output io.Writer) *Client {
	return &Client{
		conn:   conn,
		input:  bufio.NewScanner(input),
		output: output,
	}
}

// Run coordinates terminal input with server-driven termination.
func (c *Client) Run() error {
	serverDone := make(chan struct{})
	go func() {
		c.readServer()
		close(serverDone)
	}()

	c.renderBoard()
	inputDone := make(chan error, 1)
	go func() {
		inputDone <- c.processInput()
	}()

	select {
	case err := <-inputDone:
		return err
	case <-serverDone:
		return nil
	}
}
