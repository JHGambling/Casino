package server

import "jhgambling/backend/core/utils"

type GatewayClient struct {
	ID           string
	Addr         string
	IncomingChan chan []byte // Channel for messages coming from WebSocket
	OutgoingChan chan []byte // Channel for messages going to WebSocket
}

func NewGatewayClient(addr string) *GatewayClient {
	return &GatewayClient{
		ID:           utils.GenerateID(),
		Addr:         addr,
		IncomingChan: make(chan []byte, 100),
		OutgoingChan: make(chan []byte, 100),
	}
}

// Send queues a message to be sent to the WebSocket client
func (gc *GatewayClient) Send(message []byte) {
	select {
	case gc.OutgoingChan <- message:
		// Message queued successfully
	default:
		utils.Log("error", "casino::gateway", "outgoing channel full for client: ", gc.ID)
	}
}

// HandleIncomingMessage processes a message from the WebSocket
func (gc *GatewayClient) ProcessIncomingMessage(message []byte) {
	// Process the message logic here
	// For now, just log it
	utils.Log("info", "casino::gateway", "processing message from client: ", gc.ID)
	gc.Send(message)
}
