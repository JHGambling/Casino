package data

import (
	"fmt"
	"jhgambling/backend/core/data/models"
	"jhgambling/backend/core/data/tables"
	"jhgambling/backend/core/utils"
	"time"
)

// This file demonstrates how to use the table system

// Example of creating a custom game table that extends the base table
type GameTable struct {
	tables.BaseTable
}

// Custom model for the game table
type GameModel struct {
	models.UserModel        // Embedding the UserModel for demonstration
	GameID           uint   `gorm:"primaryKey"`
	Name             string `gorm:"size:255;not null"`
	Description      string `gorm:"size:1000"`
	MinimumBet       float64
	MaximumBet       float64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// NewGameTable creates a new game table
func NewGameTable() *GameTable {
	return &GameTable{
		BaseTable: tables.BaseTable{
			ID:    "games",
			Model: &GameModel{},
		},
	}
}

// Override HasPermission to implement custom permission logic
func (t *GameTable) HasPermission(userID uint, action string) bool {
	// Example permission logic
	// Only allow certain users to create/update/delete games
	if action == "read" {
		return true // Everyone can read game data
	}

	// For actions like create/update/delete, we could check if the user is an admin
	// This is just an example - you'd implement real logic here
	var user models.UserModel
	if err := t.DB.First(&user, userID).Error; err != nil {
		return false
	}

	// You could add an IsAdmin field to UserModel and check it here
	// For now, let's pretend user ID 1 is always an admin
	return user.ID == 1
}

// ExampleUsage demonstrates how to use the table system
func ExampleUsage() {
	// Create a new database instance
	db := NewDatabase()
	db.Connect()

	// Register custom tables
	userTable := tables.NewUserTable()
	gameTable := NewGameTable()

	// Register tables with the database
	db.RegisterTable(userTable)
	db.RegisterTable(gameTable)

	// Migrate all registered tables
	db.Migrate()

	// Now you can use the tables for operations
	// Create a new user
	newUser := &models.UserModel{
		Username:     "player1",
		DisplayName:  "Player One",
		PasswordHash: "hashed_password",
		JoinedAt:     time.Now(),
	}

	// Create user through the table
	err := userTable.Create(newUser)
	if err != nil {
		utils.Log("error", "example", fmt.Sprintf("Failed to create user: %v", err))
		return
	}

	utils.Log("ok", "example", fmt.Sprintf("Created user: %s", newUser.Username))

	// Create a new game
	newGame := &GameModel{
		Name:        "Blackjack",
		Description: "A card game where players try to get as close to 21 as possible.",
		MinimumBet:  5.0,
		MaximumBet:  500.0,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Create game through the table
	err = gameTable.Create(newGame)
	if err != nil {
		utils.Log("error", "example", fmt.Sprintf("Failed to create game: %v", err))
		return
	}

	utils.Log("ok", "example", fmt.Sprintf("Created game: %s", newGame.Name))

	// Find a user by ID
	userInterface, err := userTable.FindByID(1)
	if err != nil {
		utils.Log("error", "example", fmt.Sprintf("Failed to find user: %v", err))
		return
	}

	foundUser, ok := userInterface.(*models.UserModel)
	if !ok {
		utils.Log("error", "example", "Invalid user type returned")
		return
	}

	utils.Log("ok", "example", fmt.Sprintf("Found user: %s", foundUser.Username))

	// Check permissions
	hasPermission := gameTable.HasPermission(foundUser.ID, "update")
	if hasPermission {
		utils.Log("ok", "example", "User has permission to update games")
	} else {
		utils.Log("info", "example", "User does not have permission to update games")
	}
}
