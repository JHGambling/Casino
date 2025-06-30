package core

import (
	"jhgambling/backend/core/auth"
	"jhgambling/backend/core/data"
	"jhgambling/backend/core/plugins"
	"jhgambling/backend/core/server"
	"jhgambling/backend/core/game"
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
	Games 	 *game.GameManager
}

func NewCasino() *CasinoCore {
	db := data.NewDatabase()
	auth := auth.NewAuthManager()
	plugins := plugins.NewPluginManager()
	games := game.NewGameManager()

	ctx := server.GatewayContext{
		Database: db,
		Auth:     auth,
		Games: 	  games,
	}
	gateway := server.NewGateway(ctx)

	casino := CasinoCore{
		Database: db,
		Gateway:  gateway,
		Server:   server.NewServer(gateway),
		Auth:     auth,
		Plugins:  plugins,
		Games:    games,
	}

	return &casino
}

func (c *CasinoCore) Init() {
	utils.Log("info", "casino::core", "initializing...")
	c.Plugins.LoadPlugins()

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
}

func (c *CasinoCore) Start() {
	utils.Log("info", "casino::core", "starting...")

	go c.Server.Start(":9000")
	for {
		c.Gateway.Subscriptions.Update()
		time.Sleep(time.Millisecond * 10)
	}
}
