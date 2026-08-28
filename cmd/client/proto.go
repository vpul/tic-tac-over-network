package main

// message mirrors the wire format (newline-delimited JSON).
type message struct {
	Type   string    `json:"type"`
	Cell   int       `json:"cell,omitempty"`
	Symbol string    `json:"symbol,omitempty"`
	Board  [9]string `json:"board,omitempty"`
	Turn   string    `json:"turn,omitempty"`
	Result string    `json:"result,omitempty"`
	Reason string    `json:"reason,omitempty"`
}
