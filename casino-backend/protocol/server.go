package protocol

import "encoding/json"

type GamePacket struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Nonce   uint64          `json:"nonce,omitempty"`
}

type GameClient struct {
	ID string
}
