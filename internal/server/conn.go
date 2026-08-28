package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"tic-tac-over-network/internal/protocol"
)

// payload carries a decoded client request to the active session.
type payload struct {
	message protocol.Request
}

// client owns connection state. Session-specific values are assigned by the session.
type client struct {
	conn         net.Conn
	payloads     chan payload
	disconnected chan struct{}
	sessionReady chan struct{}
	session      *session
	writeMu      sync.Mutex
}

func newClient(conn net.Conn) *client {
	return &client{
		conn:         conn,
		payloads:     make(chan payload),
		disconnected: make(chan struct{}),
		sessionReady: make(chan struct{}),
	}
}

// handleConn owns the connection lifetime and forwards complete payloads to
// the matched session.
func (s *Server) handleConn(conn net.Conn) {
	connectedClient := newClient(conn)
	defer connectedClient.close()

	fmt.Printf("client connected: %s\n", connectedClient.remote())
	defer fmt.Printf("client disconnected: %s\n", connectedClient.remote())

	go connectedClient.readPayloads()
	s.waitingClients <- connectedClient

	activeSession := connectedClient.waitForSession()
	if activeSession == nil {
		return
	}
	defer activeSession.disconnect(connectedClient)

	connectedClient.processMessages(activeSession)
}

func (c *client) waitForSession() *session {
	select {
	case <-c.sessionReady:
		return c.session
	case <-c.disconnected:
		return nil
	}
}

func (c *client) processMessages(activeSession *session) {
	for {
		select {
		case incoming := <-c.payloads:
			activeSession.handle(c, incoming)
		case <-c.disconnected:
			return
		}
	}
}

func (c *client) readPayloads() {
	defer close(c.disconnected)
	decoder := json.NewDecoder(c.conn)
	for {
		var message protocol.Request
		if err := decoder.Decode(&message); err != nil {
			if !errors.Is(err, io.EOF) {
				fmt.Printf("read error from %s: %v\n", c.remote(), err)
			}
			return
		}
		c.payloads <- payload{message: message}
	}
}

func (c *client) send(message protocol.Response) {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := json.NewEncoder(c.conn).Encode(message); err != nil {
		fmt.Printf("write error to %s: %v\n", c.remote(), err)
	}
}

func (c *client) isDisconnected() bool {
	select {
	case <-c.disconnected:
		return true
	default:
		return false
	}
}

func (c *client) remote() net.Addr { return c.conn.RemoteAddr() }
func (c *client) close()           { _ = c.conn.Close() }
