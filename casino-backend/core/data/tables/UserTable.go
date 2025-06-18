package tables

import (
	"errors"
	"jhgambling/backend/core/data/models"

	"gorm.io/gorm"
)

// UserTable provides table operations for the UserModel
type UserTable struct {
	BaseTable
}

// NewUserTable creates a new user table
func NewUserTable() *UserTable {
	return &UserTable{
		BaseTable: BaseTable{
			ID:    "users",
			Model: &models.UserModel{},
		},
	}
}

// Create creates a new user
func (t *UserTable) Create(data interface{}) error {
	user, ok := data.(*models.UserModel)
	if !ok {
		return errors.New("invalid data type: expected *models.UserModel")
	}

	// Check if username already exists
	var existingUser models.UserModel
	result := t.DB.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		return errors.New("username already exists")
	} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error
	}

	return t.DB.Create(user).Error
}

// FindByID finds a user by ID
func (t *UserTable) FindByID(id interface{}) (interface{}, error) {
	var user models.UserModel
	result := t.DB.First(&user, id)
	return &user, result.Error
}

// FindByUsername finds a user by username
func (t *UserTable) FindByUsername(username string) (*models.UserModel, error) {
	var user models.UserModel
	result := t.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}

// FindAll retrieves all users with pagination
func (t *UserTable) FindAll(limit, offset int) ([]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var users []models.UserModel
	result := t.DB.Limit(limit).Offset(offset).Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert to []interface{}
	items := make([]interface{}, len(users))
	for i, user := range users {
		userCopy := user // Create a copy to avoid references to the same object
		items[i] = &userCopy
	}

	return items, nil
}

// Update updates a user
func (t *UserTable) Update(id interface{}, data interface{}) error {
	userData, ok := data.(*models.UserModel)
	if !ok {
		return errors.New("invalid data type: expected *models.UserModel")
	}

	// Don't allow changing username to one that already exists
	if userData.Username != "" {
		var existingUser models.UserModel
		result := t.DB.Where("username = ? AND id != ?", userData.Username, id).First(&existingUser)
		if result.Error == nil {
			return errors.New("username already exists")
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
	}

	return t.DB.Model(&models.UserModel{}).Where("id = ?", id).Updates(userData).Error
}
