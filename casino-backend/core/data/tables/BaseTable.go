package tables

import (
	"gorm.io/gorm"
)

// Table defines the interface that all table implementations must follow
type Table interface {
	// Basic information
	GetID() string
	GetModelType() interface{}

	// CRUD operations
	Create(data interface{}) error
	FindByID(id interface{}) (interface{}, error)
	FindAll(limit, offset int) ([]interface{}, error)
	Update(id interface{}, data interface{}) error
	Delete(id interface{}) error

	// Database interaction
	SetDB(db *gorm.DB)
	GetDB() *gorm.DB

	// Permissions
	HasPermission(userID uint, action string) bool
}

// BaseTable provides a default implementation of the Table interface
type BaseTable struct {
	ID    string
	DB    *gorm.DB
	Model interface{}
}

// GetID returns the table identifier
func (t *BaseTable) GetID() string {
	return t.ID
}

// GetModelType returns the model type for this table
func (t *BaseTable) GetModelType() interface{} {
	return t.Model
}

// SetDB sets the database connection
func (t *BaseTable) SetDB(db *gorm.DB) {
	t.DB = db
}

// GetDB returns the database connection
func (t *BaseTable) GetDB() *gorm.DB {
	return t.DB
}

// Create inserts a new record
func (t *BaseTable) Create(data interface{}) error {
	return t.DB.Create(data).Error
}

// FindByID retrieves a record by its ID
func (t *BaseTable) FindByID(id interface{}) (interface{}, error) {
	model := t.GetModelType()
	result := t.DB.First(model, id)
	return model, result.Error
}

// FindAll retrieves multiple records with pagination
func (t *BaseTable) FindAll(limit, offset int) ([]interface{}, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}

	var results []interface{}
	model := t.GetModelType()
	err := t.DB.Model(model).Limit(limit).Offset(offset).Find(&results).Error
	return results, err
}

// Update modifies an existing record
func (t *BaseTable) Update(id interface{}, data interface{}) error {
	return t.DB.Model(t.GetModelType()).Where("id = ?", id).Updates(data).Error
}

// Delete removes a record
func (t *BaseTable) Delete(id interface{}) error {
	return t.DB.Delete(t.GetModelType(), id).Error
}

// HasPermission checks if a user has permission to perform an action
// This is a default implementation that allows everything
// It should be overridden by specific tables to implement permission rules
func (t *BaseTable) HasPermission(userID uint, action string) bool {
	return true
}
