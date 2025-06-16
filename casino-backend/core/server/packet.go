package server

import "encoding/json"

type WebsocketPacket struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`
	Nonce   uint64          `json:"nonce,omitempty"`
}
