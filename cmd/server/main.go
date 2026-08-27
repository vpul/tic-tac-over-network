package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main() {
	addr := flag.String("addr", ":8080", "listen address")
	flag.Parse()

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "listen: %v\n", err)
		os.Exit(1)
	}
	defer ln.Close()
	fmt.Printf("server listening on %s\n", ln.Addr())

	go matchmaker()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Printf("accept: %v\n", err)
			continue
		}
		go handleConn(conn)
	}
}
