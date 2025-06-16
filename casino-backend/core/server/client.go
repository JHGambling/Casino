package server

import (
	"encoding/json"
	"jhgambling/backend/core/utils"
)

type GatewayClient struct {
	ID           string
	Addr         string
	IncomingChan chan []byte // Channel for messages coming from WebSocket
	OutgoingChan chan []byte // Channel for messages going to WebSocket

	handlerContext HandlerContext
}

func NewGatewayClient(addr string, ctx GatewayContext) *GatewayClient {
	client := &GatewayClient{
		ID:           utils.GenerateID(),
		Addr:         addr,
		IncomingChan: make(chan []byte, 100),
		OutgoingChan: make(chan []byte, 100),
	}

	client.handlerContext = HandlerContext{
		Client:         client,
		GatewayContext: ctx,
	}

	return client
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
	var packet WebsocketPacket
	if err := json.Unmarshal(message, &packet); err != nil {
		utils.Log("warn", "casino::gateway", "error unmarshaling data: ", err)
		return
	}

	gc.handleIncomingPacket(packet)
}

func (gc *GatewayClient) handleIncomingPacket(packet WebsocketPacket) {
	switch packet.Type {
	case "auth/register":
		var payload AuthRegisterPacket
		if gc.unmarshalPayload(packet.Payload, &payload) {
			payload.Handle(packet, &gc.handlerContext)
		}
		break
	case "auth/login":
		break
	case "auth/authenticate":
		break
	default:
		utils.Log("warn", "casino::gateway", "unknown packet type: ", packet.Type)
		return
	}
}

func (gc *GatewayClient) unmarshalPayload(payload json.RawMessage, v any) bool {
	if err := json.Unmarshal(payload, v); err != nil {
		utils.Log("warn", "casino::gateway", "error unmarshalling data: ", err)
		return false
	}
	return true
}
