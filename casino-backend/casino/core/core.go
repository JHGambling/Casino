package core

import (
	"jhgambling/backend/core/auth"
	"jhgambling/backend/core/data"
	"jhgambling/backend/core/plugins"
	"jhgambling/backend/core/server"
	"jhgambling/backend/core/utils"
)

type CasinoCore struct {
	Database *data.Database
	Server   *server.Server
	Gateway  *server.Gateway
	Auth     *auth.AuthManager
	Plugins  *plugins.PluginManager
}

func NewCasino() *CasinoCore {
	db := data.NewDatabase()
	auth := auth.NewAuthManager()
	plugins := plugins.NewPluginManager()

	ctx := server.GatewayContext{
		Database: db,
		Auth:     auth,
	}
	gateway := server.NewGateway(ctx)

	casino := CasinoCore{
		Database: db,
		Gateway:  gateway,
		Server:   server.NewServer(gateway),
		Auth:     auth,
		Plugins:  plugins,
	}

	return &casino
}

func (c *CasinoCore) Init() {
	utils.Log("info", "casino::core", "initializing...")
	c.Plugins.LoadPlugins()
	c.Database.Connect()
	c.Database.Migrate()
}

func (c *CasinoCore) Start() {
	utils.Log("info", "casino::core", "starting...")

	c.Server.Start(":9000")
}
