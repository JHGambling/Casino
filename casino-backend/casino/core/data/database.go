package data

import (
	"errors"
	"jhgambling/backend/core/data/tables"
	"jhgambling/backend/core/utils"
	"jhgambling/protocol"
	"jhgambling/protocol/models"

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

func (db *Database) Connect(path string) {
	connection, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		utils.Log("fatal", "casino::data", "Database.Connect() failed connecting to database")
		panic("failed to connect database")
	}

	db.connection = connection
	utils.Log("ok", "casino::data", "connected to db")

	db.RegisterDefaultTables()
}

func (db *Database) Migrate() {
	// Get all models from registered tables
	registeredTables := db.registry.GetAll()

	if len(registeredTables) == 0 {
		// If no tables are registered, fall back to direct model migration
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

func (db *Database) RegisterDefaultTables() {
	utils.Log("info", "casino::data", "registering default tables...")

	if err := db.RegisterTable(tables.NewUserTable()); err != nil {
		utils.Log("error", "casino::data", "error registering users table:", err)
		panic("failed to register default tables")
	}
}

// RegisterTable registers a table with the database
func (db *Database) RegisterTable(table protocol.Table) error {
	if db.connection != nil {
		table.SetDB(db.connection)
	}
	return db.registry.Register(table)
}

// GetTable retrieves a registered table by ID
func (db *Database) GetTable(tableID string) (protocol.Table, error) {
	return db.registry.Get(tableID)
}

// GetUserTable returns the user table if registered, or creates a default one
func (db *Database) GetUserTable() *tables.UserTable {
	table, err := db.registry.Get("users")
	if err != nil {
		panic("user table does not exist: " + err.Error())
	}

	userTable, ok := table.(*tables.UserTable)
	if !ok {
		panic("invalid user table")
	}

	return userTable
}

// GetTableAsUser gets a registered table with added type safety for AsUser operations
func (db *Database) GetTableAsUser(tableID string) (protocol.Table, error) {
	return db.registry.Get(tableID)
}

// PerformOperationAsUser performs a generic table operation as an authenticated user
func (db *Database) PerformOperationAsUser(authenticatedUser models.UserModel, tableID string,
	operation string, id interface{}, data interface{}) (interface{}, error) {

	utils.Log("debug", "casino::data", "[OP] user:", authenticatedUser.ID, " table:'", tableID, "' op:'", operation, "' id:", id, " data:", data)

	table, err := db.GetTableAsUser(tableID)
	if err != nil {
		utils.Log("warn", "casino::data", "[PerformOperationAsUser] error getting table:", err)
		return nil, err
	}

	switch operation {
	case "create":
		return nil, table.CreateAsUser(authenticatedUser, data)
	case "findByID":
		return table.FindByIDAsUser(authenticatedUser, id)
	case "findAll":
		limit, ok := id.(int)
		if !ok {
			limit = 10
		}
		offset, ok := data.(int)
		if !ok {
			offset = 0
		}
		return table.FindAllAsUser(authenticatedUser, limit, offset)
	case "update":
		return nil, table.UpdateAsUser(authenticatedUser, id, data)
	case "delete":
		return nil, table.DeleteAsUser(authenticatedUser, id)
	default:
		utils.Log("warn", "casino::data", "[PerformOperationAsUser] unkown operation \""+operation+"\" by user", authenticatedUser.ID)
		return nil, errors.New("unknown operation")
	}
}
