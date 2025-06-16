package data

import (
	"jhgambling/backend/core/data/models"
	"jhgambling/backend/core/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	connection *gorm.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) Connect() {
	connection, err := gorm.Open(sqlite.Open("casino.db"), &gorm.Config{})
	if err != nil {
		utils.Log("fatal", "casino::data", "Database.Connect() failed connecting to database")
		panic("failed to connect database")
	}

	db.connection = connection
	utils.Log("ok", "casino::data", "connected to db")
}

func (db *Database) Migrate() {
	db.connection.AutoMigrate(&models.UserModel{})

	utils.Log("ok", "casino::data", "migrated all models")
}
