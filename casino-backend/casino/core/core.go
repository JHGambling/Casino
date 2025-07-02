package core

import (
	"jhgambling/backend/core/auth"
	"jhgambling/backend/core/data"
	"jhgambling/backend/core/game"
	"jhgambling/backend/core/plugins"
	"jhgambling/backend/core/server"
	"jhgambling/backend/core/utils"
	"os"
	"time"
)

type CasinoCore struct {
	Database *data.Database
	Server   *server.Server
	Gateway  *server.Gateway
	Auth     *auth.AuthManager
	Plugins  *plugins.PluginManager
	Games    *game.GameManager

	Adapter *CasinoPluginAdapter
}

func NewCasino() *CasinoCore {
	db := data.NewDatabase()
	auth := auth.NewAuthManager()
	plugins := plugins.NewPluginManager()
	games := game.NewGameManager()

	ctx := server.GatewayContext{
		Database: db,
		Auth:     auth,
		Games:    games,
	}
	gateway := server.NewGateway(ctx)

	casino := &CasinoCore{
		Database: db,
		Gateway:  gateway,
		Server:   server.NewServer(gateway),
		Auth:     auth,
		Plugins:  plugins,
		Games:    games,
	}

	adapter := NewCasinoPluginAdapter(casino)
	casino.Adapter = adapter

	return casino
}

func (c *CasinoCore) Init() {
	utils.Log("info", "casino::core", "initializing...")

	// plugins
	c.Plugins.LoadPlugins()

	// Database
	env := os.Getenv("ENV")
	dbPath := ""
	if env == "production" {
		dbPath = "/data/casino.db"
	} else {
		dbPath = "../casino.db"
	}
	c.Database.Connect(dbPath)
	c.Database.Migrate()
	c.Database.SetSubscriptionChannel(&c.Gateway.Subscriptions.ChangedRecordsChannel)

	// Game integration
	c.Games.SetAdapter(c.Adapter)
	c.registerGameProviders()
}

func (c *CasinoCore) Start() {
	utils.Log("info", "casino::core", "starting...")

	go c.Server.Start(":9000")
	for {
		c.Gateway.Subscriptions.Update()
		time.Sleep(time.Millisecond * 10)
	}
}

func (c *CasinoCore) registerGameProviders() {
	for _, p := range c.Plugins.GameProviders {
		c.Games.RegisterProvider(p)
	}
}
