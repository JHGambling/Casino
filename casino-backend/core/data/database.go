package data

import (
	"errors"
	"jhgambling/backend/core/data/models"
	"jhgambling/backend/core/data/tables"
	"jhgambling/backend/core/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Database struct {
	connection *gorm.DB
	registry   *tables.TableRegistry
}

func NewDatabase() *Database {
	return &Database{
		registry: tables.NewTableRegistry(),
	}
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
	// Get all models from registered tables
	registeredTables := db.registry.GetAll()

	if len(registeredTables) == 0 {
		// If no tables are registered, fall back to direct model migration
		db.connection.AutoMigrate(&models.UserModel{})
	} else {
		// Migrate models from registered tables
		for _, table := range registeredTables {
			db.connection.AutoMigrate(table.GetModelType())

			// Set the DB connection for each table
			table.SetDB(db.connection)
		}
	}

	utils.Log("ok", "casino::data", "migrated all models")
}

// RegisterTable registers a table with the database
func (db *Database) RegisterTable(table tables.Table) error {
	if db.connection != nil {
		table.SetDB(db.connection)
	}
	return db.registry.Register(table)
}

// GetTable retrieves a registered table by ID
func (db *Database) GetTable(tableID string) (tables.Table, error) {
	return db.registry.Get(tableID)
}

// GetUserTable returns the user table if registered, or creates a default one
func (db *Database) GetUserTable() (*tables.UserTable, error) {
	table, err := db.registry.Get("users")
	if err != nil {
		// Create a default user table if not registered
		userTable := tables.NewUserTable()
		userTable.SetDB(db.connection)
		err = db.registry.Register(userTable)
		if err != nil {
			return nil, err
		}
		return userTable, nil
	}

	userTable, ok := table.(*tables.UserTable)
	if !ok {
		return nil, errors.New("registered users table is not of type UserTable")
	}

	return userTable, nil
}

// Backward compatibility methods that use the user table

func (db *Database) GetUserByID(userID uint) (models.UserModel, bool, error) {
	userTable, err := db.GetUserTable()
	if err != nil {
		return models.UserModel{}, false, err
	}

	user, err := userTable.FindByID(userID)
	if err != nil {
		return models.UserModel{}, false, err
	}

	userModel, ok := user.(*models.UserModel)
	if !ok {
		return models.UserModel{}, false, errors.New("invalid user model type")
	}

	return *userModel, true, nil
}

func (db *Database) GetUserByUsername(username string) (models.UserModel, bool, error) {
	userTable, err := db.GetUserTable()
	if err != nil {
		return models.UserModel{}, false, err
	}

	user, err := userTable.FindByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.UserModel{}, false, nil
		}
		return models.UserModel{}, false, err
	}

	return *user, true, nil
}

func (db *Database) CreateUser(user *models.UserModel) error {
	userTable, err := db.GetUserTable()
	if err != nil {
		return err
	}

	return userTable.Create(user)
}
