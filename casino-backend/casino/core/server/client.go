package server

import (
	"encoding/json"
	"jhgambling/backend/core/utils"
	"jhgambling/protocol"
	"time"
)

type GatewayClient struct {
	ID           string
	Addr         string
	IncomingChan chan []byte // Channel for messages coming from WebSocket
	OutgoingChan chan []byte // Channel for messages going to WebSocket

	handlerContext HandlerContext

	isAuthenticated         bool
	authenticatedAs         uint
	authenticationExpriesAt time.Time

	Subscriptions []DBSubscription
}

func NewGatewayClient(addr string, ctx GatewayContext) *GatewayClient {
	client := &GatewayClient{
		ID:           utils.GenerateID(),
		Addr:         addr,
		IncomingChan: make(chan []byte, 100),
		OutgoingChan: make(chan []byte, 100),

		Subscriptions: []DBSubscription{},
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
		var payload AuthLoginPacket
		if gc.unmarshalPayload(packet.Payload, &payload) {
			payload.Handle(packet, &gc.handlerContext)
		}
		break
	case "auth/authenticate":
		var payload AuthAuthenticatePacket
		if gc.unmarshalPayload(packet.Payload, &payload) {
			payload.Handle(packet, &gc.handlerContext)
		}
		break
	case "auth/does_user_exist":
		var payload DoesUserExistPacket
		if gc.unmarshalPayload(packet.Payload, &payload) {
			payload.Handle(packet, &gc.handlerContext)
		}
		break
	case "db/op":
		var payload DatabaseOperationPacket
		if gc.unmarshalPayload(packet.Payload, &payload) {
			payload.Handle(packet, &gc.handlerContext)
		}
		break
	case "db/sub":
		var payload DatabaseSubscribePacket
		if gc.unmarshalPayload(packet.Payload, &payload) {
			payload.Handle(packet, &gc.handlerContext)
		}
		break
	case "ping":
		// Simply respond with a pong to keep the connection alive
		if res, err := BuildPacket("pong", map[string]interface{}{}, packet.Nonce); err == nil {
			gc.Send(res)
		}
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

func (gc *GatewayClient) IsAuthenticated() bool {
	return gc.isAuthenticated && time.Now().Unix() < gc.authenticationExpriesAt.Unix()
}

func (gc *GatewayClient) Authenticate(userID uint, expiresAt time.Time) {
	gc.isAuthenticated = true
	gc.authenticatedAs = userID
	gc.authenticationExpriesAt = expiresAt
}

func (gc *GatewayClient) RevokeAuthentication() {
	gc.isAuthenticated = false
	gc.authenticatedAs = 0
	gc.authenticationExpriesAt = time.UnixMicro(0)
}

func (gc *GatewayClient) SendUnauthorizedPacket(nonce uint64) {
	if res, err := BuildPacket("res",
		ResponsePacket{
			Success: false,
			Status:  "unauthorized",
			Message: "You have to be authorized to interact",
		},
		nonce); err == nil {
		gc.Send(res)
	}
}

func (gc *GatewayClient) SendSubscriptionUpdatePacket(record protocol.SubChangedRecord) {
	if res, err := BuildPacket("db/sub:update",
		DatabaseSubUpdatePacket{
			Operation:  record.Operation,
			TableID:    record.TableID,
			ResourceID: record.ResourceID,
			Data:       record.Record,
		},
		0); err == nil {
		gc.Send(res)
	}
}
