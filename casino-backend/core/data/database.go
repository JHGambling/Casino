package data

import (
	"errors"
	"jhgambling/backend/core/data/models"
	"jhgambling/backend/core/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	connection *gorm.DB
}

func NewDatabase() *Database {
	return &Database{}
}

func (db *Database) Connect() {
	connection, err := gorm.Open(sqlite.Open("casino.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
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

func (db *Database) GetUserByID(userID uint) (models.UserModel, bool, error) {
	var user models.UserModel
	result := db.connection.First(&user, userID)

	return user, !errors.Is(result.Error, gorm.ErrRecordNotFound), result.Error
}

func (db *Database) GetUserByUsername(username string) (models.UserModel, bool, error) {
	var user models.UserModel
	result := db.connection.Where("username = ?", username).First(&user)

	return user, !errors.Is(result.Error, gorm.ErrRecordNotFound), result.Error
}

func (db *Database) CreateUser(user *models.UserModel) error {
	result := db.connection.Create(user)
	return result.Error
}
