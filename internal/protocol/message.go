package protocol

// Request is a client command encoded as newline-delimited JSON.
type Request struct {
	Type string `json:"type"` // "move"
	Cell int    `json:"cell,omitempty"`
}

// Response is a server event encoded as newline-delimited JSON.
type Response struct {
	Type   string    `json:"type"` // "waiting", "game_start", "state", "game_over", or "error"
	Symbol string    `json:"symbol,omitempty"`
	Board  [9]string `json:"board,omitempty"`
	Turn   string    `json:"turn,omitempty"`
	Result string    `json:"result,omitempty"`
	Reason string    `json:"reason,omitempty"`
}
