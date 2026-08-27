package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
)

// symbol is written by readServer (on game_start) and read by the input
// loop on every move.
var (
	symbolMu sync.Mutex
	symbol   = "?" // until the server assigns one at pairing
)

func currentSymbol() string {
	symbolMu.Lock()
	defer symbolMu.Unlock()
	return symbol
}

func setSymbol(s string) {
	symbolMu.Lock()
	defer symbolMu.Unlock()
	symbol = s
}

// readServer parses server messages until the conn closes.
func readServer(conn net.Conn) {
	dec := json.NewDecoder(conn)
	for {
		var m message
		if err := dec.Decode(&m); err != nil {
			if errors.Is(err, io.EOF) {
				fmt.Println("connection closed by server")
			} else {
				fmt.Fprintf(os.Stderr, "read error: %v\n", err)
			}
			os.Exit(0)
		}
		switch m.Type {
		case "waiting":
			fmt.Println("waiting for an opponent...")
		case "game_start":
			setSymbol(m.Symbol)
			fmt.Printf("paired! you are %s\n", m.Symbol)
		default:
			fmt.Printf("[server] %+v\n", m)
		}
	}
}

func main() {
	addr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	out := json.NewEncoder(conn)
	if err := out.Encode(message{Type: "hello"}); err != nil {
		fmt.Fprintf(os.Stderr, "send hello: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("connected to %s\n", *addr)

	go readServer(conn)

	// Local preview only. The server owns the real board once it validates;
	// this display is a mirror, never an authority.
	board := [9]string{"1", "2", "3", "4", "5", "6", "7", "8", "9"}
	render(board)

	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		line := strings.TrimSpace(in.Text())
		switch {
		case line == "":
			continue
		case line == "q" || line == "quit":
			return
		default:
			n, err := strconv.Atoi(line)
			if err != nil || n < 1 || n > 9 {
				fmt.Println("? enter a cell 1-9, or q to quit")
				continue
			}
			sym := currentSymbol()
			if sym == "?" {
				fmt.Println("not paired yet — wait for an opponent")
				continue
			}
			mark := strconv.Itoa(n)
			if board[n-1] != mark {
				fmt.Printf("(local preview) cell %d already marked %s — not sent\n", n, board[n-1])
				continue
			}
			board[n-1] = sym
			render(board)
			if err := out.Encode(message{Type: "move", Cell: n}); err != nil {
				fmt.Fprintf(os.Stderr, "send move: %v\n", err)
				os.Exit(1)
			}
		}
	}
}
