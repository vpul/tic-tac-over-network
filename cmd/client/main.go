package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
)

func main() {
	addr := flag.String("addr", "localhost:8080", "server address")
	flag.Parse()

	conn, err := net.Dial("tcp", *addr)
	if err != nil {
		log.Fatalf("dial: %v", err)
	}
	defer conn.Close()

	if _, err := fmt.Fprintln(conn, "hello"); err != nil {
		log.Fatalf("send: %v", err)
	}
	log.Println("sent: hello")

	ack, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Fatalf("read ack: %v", err)
	}
	fmt.Printf("received: %s", ack)
}
