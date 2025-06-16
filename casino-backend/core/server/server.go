package server

import (
	"encoding/json"
	"jhgambling/backend/core/utils"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Addr        string
	SendChannel chan []byte
	RecvChannel chan []byte
	GatewayID   string // Link to the gateway client
}

type Server struct {
	clients    map[string]*Client
	gateway    *Gateway
	mu         sync.Mutex
	upgrader   websocket.Upgrader
	httpServer *http.Server
}

func NewServer(gateway *Gateway) *Server {
	return &Server{
		clients: make(map[string]*Client),
		gateway: gateway,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

func (s *Server) Start(addr string) error {
	http.HandleFunc("/ws", s.handleWebSocket)
	http.HandleFunc("/api", s.handleAPI)

	s.httpServer = &http.Server{Addr: addr}
	utils.Log("info", "casino::server", "starting server on ", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) handleAPI(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"message": "test"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		utils.Log("error", "casino::server", "error upgrading websocket:", err)
		return
	}
	defer conn.Close()

	// Create a new Gateway client
	gatewayClient := NewGatewayClient(r.RemoteAddr)

	// Add client to gateway
	s.gateway.AddClient(gatewayClient)

	// Create WebSocket client linked to the gateway client
	client := &Client{
		Addr:        r.RemoteAddr,
		SendChannel: make(chan []byte, 100),
		RecvChannel: make(chan []byte, 100),
		GatewayID:   gatewayClient.ID,
	}

	s.mu.Lock()
	s.clients[client.Addr] = client
	s.mu.Unlock()

	// Start goroutines to handle communication between WebSocket and Gateway
	go s.handleClientToGateway(client, gatewayClient)
	go s.handleGatewayToClient(client, gatewayClient, conn)

	// Read messages from WebSocket
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			utils.Log("error", "casino::server", "error reading websocket:", err)
			break
		}
		client.RecvChannel <- message
	}

	// Clean up when connection closes
	s.mu.Lock()
	delete(s.clients, client.Addr)
	s.mu.Unlock()
	s.gateway.RemoveClient(gatewayClient.ID)
	utils.Log("info", "casino::server", "client disconnected:", r.RemoteAddr)
}

// Handle messages from WebSocket client to Gateway
func (s *Server) handleClientToGateway(client *Client, gatewayClient *GatewayClient) {
	for {
		select {
		case msg, ok := <-client.RecvChannel:
			if !ok {
				return // Channel closed
			}

			// Forward message to Gateway client's incoming channel
			gatewayClient.IncomingChan <- msg

			// Process the message in the Gateway client
			//gatewayClient.ProcessIncomingMessage(msg)
		}
	}
}

// Handle messages from Gateway to WebSocket client
func (s *Server) handleGatewayToClient(client *Client, gatewayClient *GatewayClient, conn *websocket.Conn) {
	for {
		select {
		case msg, ok := <-gatewayClient.OutgoingChan:
			if !ok {
				return // Channel closed
			}

			// Send message to WebSocket
			err := conn.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				utils.Log("error", "casino::server", "error writing to websocket:", err)
				return
			}
		}
	}
}
