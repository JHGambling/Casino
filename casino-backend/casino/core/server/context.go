package server

import (
	"jhgambling/backend/core/auth"
	"jhgambling/backend/core/data"
	"jhgambling/backend/core/game"
)

type GatewayContext struct {
	Database *data.Database
	Auth     *auth.AuthManager
	Gateway  *Gateway
	Games    *game.GameManager
}

type HandlerContext struct {
	Client *GatewayClient

	GatewayContext
}
