package protocol

import (
	"jhgambling/protocol/models"

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

	// User Operations
	CreateAsUser(user models.UserModel, data interface{}) error
	FindByIDAsUser(user models.UserModel, id interface{}) (interface{}, error)
	FindAllAsUser(user models.UserModel, limit, offset int) ([]interface{}, error)
	UpdateAsUser(user models.UserModel, id interface{}, data interface{}) error
	DeleteAsUser(user models.UserModel, id interface{}) error

	// Database interaction
	SetDB(db *gorm.DB)
	GetDB() *gorm.DB

	SetSubscriptionChannel(channel *chan SubChangedRecord)
	PushRecordChange(operation string, id interface{}, data interface{})
	CanViewChangedRecord(user models.UserModel, record SubChangedRecord) bool

	// Single-time functions
	Repair()
}

// BaseTable provides a default implementation of the Table interface
type BaseTable struct {
	ID                  string
	DB                  *gorm.DB
	SubscriptionChannel *chan SubChangedRecord
	Model               interface{}
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

func (t *BaseTable) SetSubscriptionChannel(channel *chan SubChangedRecord) {
	t.SubscriptionChannel = channel
}

// This function can be used to repair models when the database starts up
func (t *BaseTable) Repair() {

}

// Create inserts a new record
func (t *BaseTable) Create(data interface{}) error {
	err := t.DB.Create(data).Error
	if err == nil {
		t.PushRecordChange("create", nil, data)
	}

	return err
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
	// Perform the update
	err := t.DB.Model(t.GetModelType()).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return err
	}

	// Create a new instance of the model to hold the updated row
	updated := t.GetModelType()

	// Fetch the updated record
	err = t.DB.First(updated, "id = ?", id).Error
	if err != nil {
		return err
	}

	// Pass the updated record to PushRecordChange
	t.PushRecordChange("update", id, updated)

	return nil
}

// Delete removes a record
func (t *BaseTable) Delete(id interface{}) error {
	err := t.DB.Delete(t.GetModelType(), id).Error
	if err == nil {
		t.PushRecordChange("delete", id, nil)
	}

	return err
}

// CreateAsUser creates a new record with user permission check
func (t *BaseTable) CreateAsUser(user models.UserModel, data interface{}) error {
	// Add Permission check
	return t.Create(data)
}

// FindByIDAsUser retrieves a record by ID with user permission check
func (t *BaseTable) FindByIDAsUser(user models.UserModel, id interface{}) (interface{}, error) {
	// Add Permission check
	return t.FindByID(id)
}

// FindAllAsUser retrieves multiple records with pagination and user permission check
func (t *BaseTable) FindAllAsUser(user models.UserModel, limit, offset int) ([]interface{}, error) {
	// Add Permission check
	return t.FindAll(limit, offset)
}

// UpdateAsUser modifies an existing record with user permission check
func (t *BaseTable) UpdateAsUser(user models.UserModel, id interface{}, data interface{}) error {
	// Add Permission check
	return t.Update(id, data)
}

// DeleteAsUser removes a record with user permission check
func (t *BaseTable) DeleteAsUser(user models.UserModel, id interface{}) error {
	// Add Permission check
	return t.Delete(id)
}

func (t *BaseTable) PushRecordChange(operation string, id interface{}, data interface{}) {
	*t.SubscriptionChannel <- SubChangedRecord{
		Operation:  operation,
		TableID:    t.ID,
		ResourceID: id,
		Record:     data,
	}
}

func (t *BaseTable) CanViewChangedRecord(user models.UserModel, record SubChangedRecord) bool {
	// Add Permission check
	return true
}
