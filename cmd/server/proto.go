package main

// message mirrors the client's wire type (cmd/client/proto.go).
// Next increment (update/game_over) moves both copies to internal/proto.
type message struct {
	Type   string `json:"type"`
	Cell   int    `json:"cell,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}
