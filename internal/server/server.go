package server

import (
	"fmt"
	"net"
)

type Server struct {
	waitingClients chan *client
}

func New() *Server {
	return &Server{waitingClients: make(chan *client, 1)}
}

func (s *Server) Serve(listener net.Listener) error {
	fmt.Printf("server listening on %s\n", listener.Addr())
	go s.runLobby()

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("accept: %w", err)
		}
		go s.handleConn(conn)
	}
}
