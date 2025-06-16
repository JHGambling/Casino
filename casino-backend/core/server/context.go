package server

import (
	"jhgambling/backend/core/auth"
	"jhgambling/backend/core/data"
)

type GatewayContext struct {
	Database *data.Database
	Auth     *auth.AuthManager
	Gateway  *Gateway
}

type HandlerContext struct {
	Client *GatewayClient

	GatewayContext
}
