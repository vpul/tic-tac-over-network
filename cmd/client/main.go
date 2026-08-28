package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	clientpkg "tic-tac-over-network/internal/client"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "dial: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Printf("connected to %s\n", *addr)
	if err := clientpkg.New(conn, os.Stdin, os.Stdout).Run(); err != nil {
		fmt.Fprintf(os.Stderr, "client: %v\n", err)
		os.Exit(1)
	}
}
