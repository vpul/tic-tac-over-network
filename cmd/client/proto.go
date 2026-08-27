package main

// message mirrors the wire format (newline-delimited JSON).
type message struct {
	Type   string `json:"type"`
	Cell   int    `json:"cell,omitempty"`
	Symbol string `json:"symbol,omitempty"`
}
