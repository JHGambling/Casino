package core

import (
	"jhgambling/backend/core/data"
	"jhgambling/backend/core/server"
	"jhgambling/backend/core/utils"
)

type CasinoCore struct {
	Database *data.Database
	Server   *server.Server
	Gateway  *server.Gateway
}

func NewCasino() *CasinoCore {
	gateway := server.NewGateway()

	casino := CasinoCore{
		Database: data.NewDatabase(),
		Gateway:  gateway,
		Server:   server.NewServer(gateway),
	}

	return &casino
}

func (c *CasinoCore) Init() {
	utils.Log("info", "casino::core", "initializing...")
	c.Database.Connect()
	c.Database.Migrate()
}

func (c *CasinoCore) Start() {
	utils.Log("info", "casino::core", "starting...")

	c.Server.Start(":9000")
}
