package server

import (
	"encoding/json"
	"jhgambling/backend/core/utils"
)

// BuildPacket creates a WebsocketPacket and serializes it into a sendable JSON format.
func BuildPacket(packetType string, payload interface{}, nonce uint64) ([]byte, error) {
	// Serialize the payload into JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		utils.Log("error", "casino::server", "error marshaling payload:", err)
		return nil, err
	}

	// Create the WebsocketPacket
	wsPacket := WebsocketPacket{
		Type:    packetType,
		Payload: payloadBytes,
		Nonce:   nonce,
	}

	// Serialize the WebsocketPacket into JSON
	packetBytes, err := json.Marshal(wsPacket)
	if err != nil {
		utils.Log("error", "casino::server", "error marshaling websocket packet:", err)
		return nil, err
	}

	return packetBytes, nil
}
