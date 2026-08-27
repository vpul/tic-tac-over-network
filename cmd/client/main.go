package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

// message mirrors internal/proto's eventual wire format (newline-delimited JSON).
type message struct {
	Type string `json:"type"`
	Cell int    `json:"cell,omitempty"`
}

func render(board [9]string) {
	fmt.Println()
	for r := 0; r < 3; r++ {
		i := r * 3
		fmt.Printf(" %s | %s | %s \n", board[i], board[i+1], board[i+2])
		if r < 2 {
			fmt.Println("---+---+---")
		}
	}
	fmt.Println("\nplay a cell (1-9) or q to quit:")
}

func main() {
	addr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	out := json.NewEncoder(conn)
	if err := out.Encode(message{Type: "hello"}); err != nil {
		log.Fatalf("send hello: %v", err)
	}
	log.Printf("connected to %s", *addr)

	// server -> screen: print everything until the server closes the conn.
	go func() {
		scanner := bufio.NewScanner(conn)
		for scanner.Scan() {
			fmt.Printf("[server] %s\n", scanner.Text())
		}
		log.Println("connection closed by server")
		os.Exit(0)
	}()

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
			mark := strconv.Itoa(n)
			if board[n-1] != mark {
				fmt.Printf("(local preview) cell %d already marked %s — not sent\n", n, board[n-1])
				continue
			}
			board[n-1] = "X"
			render(board)
			if err := out.Encode(message{Type: "move", Cell: n}); err != nil {
				log.Fatalf("send move: %v", err)
			}
		}
	}
}
