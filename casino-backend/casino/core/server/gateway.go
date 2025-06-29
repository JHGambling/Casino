package server

import (
	"errors"
	"jhgambling/backend/core/utils"
	"sync"
)

type Gateway struct {
	Clients map[string]*GatewayClient
	mu      sync.Mutex

	Subscriptions *SubscriptionManager

	ctx GatewayContext
}

func NewGateway(ctx GatewayContext) *Gateway {
	gw := &Gateway{
		Clients: make(map[string]*GatewayClient),
		ctx:     ctx,
	}

	gw.Subscriptions = NewSubscriptionsManager(gw)
	gw.ctx.Gateway = gw

	return gw
}

// AddClient adds a new GatewayClient to the Gateway
func (g *Gateway) AddClient(client *GatewayClient) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Clients[client.ID] = client
	g.StartClientHandler(client) // Start handling messages for this client
	utils.Log("info", "casino::gateway", "client added:", client.ID)
}

func (g *Gateway) RemoveClient(clientID string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if _, exists := g.Clients[clientID]; exists {
		delete(g.Clients, clientID)
		utils.Log("info", "casino::gateway", ">> client removed: ", clientID)
	}
}

// Broadcast sends a message to all connected clients
func (g *Gateway) Broadcast(message []byte) {
	g.mu.Lock()
	defer g.mu.Unlock()
	for _, client := range g.Clients {
		client.Send(message)
	}
}

// SendToClient sends a message to a specific client by ID
func (g *Gateway) SendToClient(clientID string, message []byte) error {
	g.mu.Lock()
	client, exists := g.Clients[clientID]
	g.mu.Unlock()

	if !exists {
		return errors.New("client not found")
	}

	client.Send(message)
	return nil
}

func (g *Gateway) StartClientHandler(client *GatewayClient) {
	go func() {
		for {
			select {
			case msg, ok := <-client.IncomingChan:
				if !ok {
					// Channel closed
					return
				}
				// Process the message
				client.ProcessIncomingMessage(msg)
			}
		}
	}()
}
