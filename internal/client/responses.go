package client

import (
	"encoding/json"
	"errors"
	"io"

	"tic-tac-over-network/internal/game"
	"tic-tac-over-network/internal/protocol"
)

func (c *Client) readServer() {
	decoder := json.NewDecoder(c.conn)
	for {
		var response protocol.Response
		if err := decoder.Decode(&response); err != nil {
			if errors.Is(err, io.EOF) {
				c.println("connection closed by server")
			} else {
				c.printf("read error: %v\n", err)
			}
			return
		}
		if c.handleResponse(response) {
			return
		}
	}
}

func (c *Client) handleResponse(response protocol.Response) bool {
	switch response.Type {
	case "waiting":
		c.println("waiting for an opponent...")
	case "game_start":
		c.stateMu.Lock()
		c.symbol = response.Symbol
		c.stateMu.Unlock()
		c.updateState(response.Board, response.Turn)
		c.printf("paired! you are %s\n", response.Symbol)
		c.renderBoard()
	case "state":
		c.updateState(response.Board, response.Turn)
		c.renderBoard()
	case "game_over":
		c.updateState(response.Board, response.Turn)
		c.renderBoard()
		c.printResult(response.Result)
		return true
	case "opponent_left":
		c.println("opponent left the game")
		return true
	case "error":
		c.printf("move rejected: %s\n", response.Reason)
	default:
		c.printf("[server] %+v\n", response)
	}
	return false
}

func (c *Client) updateState(board game.Board, turn string) {
	c.stateMu.Lock()
	defer c.stateMu.Unlock()
	c.board = board
	c.turn = turn
}
func (c *Client) printResult(result string) {
	c.stateMu.Lock()
	symbol := c.symbol
	c.stateMu.Unlock()
	switch {
	case result == "draw":
		c.println("game over: draw")
	case result == symbol:
		c.println("game over: you win")
	default:
		c.printf("game over: %s wins\n", result)
	}
}
