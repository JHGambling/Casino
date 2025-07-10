package tables

import (
	"errors"
	"jhgambling/backend/core/utils"
	"jhgambling/protocol"
	"jhgambling/protocol/models"

	"gorm.io/gorm"
)

// SafeUserModel represents a user model with sensitive information removed
type SafeUserModel struct {
	ID          uint
	Username    string
	DisplayName string
	JoinedAt    string
	IsAdmin     bool
	Wallet      models.WalletModel
}

// UserTable provides table operations for the UserModel
type UserTable struct {
	protocol.BaseTable
}

// NewUserTable creates a new user table
func NewUserTable() *UserTable {
	return &UserTable{
		BaseTable: protocol.BaseTable{
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
	_, err := t.FindByUsername(user.Username)
	if err == nil {
		return errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	err = t.DB.Create(user).Error
	if err == nil {
		t.PushRecordChange("create", user.ID, toSafeUser(user))
	}

	return err
}

// FindByID finds a user by ID
func (t *UserTable) FindByID(id interface{}) (interface{}, error) {
	var user models.UserModel
	result := t.DB.Preload("Wallet").First(&user, id)
	return &user, result.Error
}

// FindByUsername finds a user by username
func (t *UserTable) FindByUsername(username string) (*models.UserModel, error) {
	var user models.UserModel
	result := t.DB.Where("username = ?", username).Preload("Wallet").First(&user)
	return &user, result.Error
}

// FindAll retrieves all users with pagination
func (t *UserTable) FindAll(limit, offset int) ([]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var users []models.UserModel
	result := t.DB.Limit(limit).Offset(offset).Preload("Wallet").Find(&users)
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

// CreateAsUser implements user-based permission check for creating a user
func (t *UserTable) CreateAsUser(user models.UserModel, data interface{}) error {
	return t.Create(data)
}

// toSafeUser converts a UserModel to a SafeUserModel by removing sensitive information
func toSafeUser(user *models.UserModel) *SafeUserModel {
	return &SafeUserModel{
		ID:          user.ID,
		Username:    user.Username,
		DisplayName: user.DisplayName,
		JoinedAt:    user.JoinedAt.Format("2006-01-02 15:04:05"),
		IsAdmin:     user.IsAdmin,
		Wallet:      user.Wallet,
	}
}

// FindByIDAsUser retrieves a user by ID with permission check and removes sensitive data
func (t *UserTable) FindByIDAsUser(user models.UserModel, id interface{}) (interface{}, error) {
	foundUser, err := t.FindByID(id)
	if err != nil {
		return nil, err
	}

	userModel, ok := foundUser.(*models.UserModel)
	if !ok {
		return nil, errors.New("invalid user model type")
	}

	return toSafeUser(userModel), nil
}

// FindAllAsUser retrieves all users with permission check and removes sensitive data
func (t *UserTable) FindAllAsUser(user models.UserModel, limit, offset int) ([]interface{}, error) {
	users, err := t.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	// Convert to safe user models
	safeUsers := make([]interface{}, len(users))
	for i, u := range users {
		userModel, ok := u.(*models.UserModel)
		if !ok {
			return nil, errors.New("invalid user model type")
		}
		safeUsers[i] = toSafeUser(userModel)
	}

	return safeUsers, nil
}

// UpdateAsUser modifies a user with permission check
func (t *UserTable) UpdateAsUser(user models.UserModel, id interface{}, data interface{}) error {
	// Users can update their own data
	userID, ok := id.(uint)
	if !ok {
		return errors.New("invalid ID format")
	}

	if user.ID != userID && !user.IsAdmin {
		return errors.New("permission denied: you can only update your own user data")
	}

	// If updating the user's password, ensure we have a new hash
	userData, ok := data.(*models.UserModel)
	if !ok {
		return errors.New("invalid data type: expected *models.UserModel")
	}

	// If regular user is trying to change their username, apply special validation
	if userData.Username != "" && userData.Username != user.Username {
		// Check if new username already exists
		_, err := t.FindByUsername(userData.Username)
		if err == nil {
			return errors.New("username already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	return t.Update(id, userData)
}

// DeleteAsUser removes a user with permission check
func (t *UserTable) DeleteAsUser(user models.UserModel, id interface{}) error {
	// Only admins can delete users
	if !user.IsAdmin {
		return errors.New("permission denied: only admins can delete users")
	}

	return t.Delete(id)
}

// Update updates a user
func (t *UserTable) Update(id interface{}, data interface{}) error {
	userData, ok := data.(*models.UserModel)
	if !ok {
		return errors.New("invalid data type: expected *models.UserModel")
	}

	// Don't allow changing username to one that already exists
	if userData.Username != "" {
		existing, err := t.FindByUsername(userData.Username)
		if err == nil && existing.ID != id {
			return errors.New("username already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
	}

	// Perform the update
	err := t.DB.Model(&models.UserModel{}).Where("id = ?", id).Updates(userData).Error
	if err != nil {
		return err
	}

	// Fetch the updated row
	var updatedUser models.UserModel
	err = t.DB.First(&updatedUser, "id = ?", id).Error
	if err != nil {
		return err
	}

	// Push the change with the updated record
	t.PushRecordChange("update", id, &updatedUser)

	return nil
}

func (t *UserTable) Repair() {
	t.repair_addWallets()
	t.repair_addStartingBonus()
}

func (t *UserTable) repair_addWallets() {
	allUsers, err := t.FindAll(10000, 0)
	if err != nil {
		utils.Log("error", "casino::data", "[UserTable] [Repair] repair_addWallets failed with an error while getting all users:", err)
		return
	}

	for _, u := range allUsers {
		userData, ok := u.(*models.UserModel)
		if !ok {
			continue
		}

		// Check if user has a wallet
		if userData.Wallet.ID == 0 {
			// Create new wallet for the user
			wallet := &models.WalletModel{
				UserID:                userData.ID,
				NetworthCents:         0,
				ReceivedStartingBonus: false,
			}

			// Save the wallet directly
			err := t.DB.Create(wallet).Error
			if err != nil {
				utils.Log("error", "casino::data", "[UserTable] [Repair] failed to create wallet for user:", userData.ID, "error:", err)
				continue
			}

			utils.Log("ok", "casino::data", "[UserTable] [Repair] created wallet for user:", userData.ID)
		}
	}
}

// repair_addStartingBonus gives a starting bonus of $1000 to users who haven't received it yet
func (t *UserTable) repair_addStartingBonus() {
	utils.Log("info", "casino::data", "[UserTable] [Repair] checking for users who need starting bonus...")

	// Get all wallets that haven't received the starting bonus
	var wallets []models.WalletModel
	result := t.DB.Where("received_starting_bonus = ?", false).Find(&wallets)
	if result.Error != nil {
		utils.Log("error", "casino::data", "[UserTable] [Repair] failed to get wallets:", result.Error)
		return
	}

	// Process each wallet
	for _, wallet := range wallets {
		// Add $1000 in cents (100000 cents)
		wallet.NetworthCents += 100000
		wallet.ReceivedStartingBonus = true

		// Update the wallet
		err := t.DB.Save(&wallet).Error
		if err != nil {
			utils.Log("error", "casino::data", "[UserTable] [Repair] failed to update wallet for user:", wallet.UserID, "error:", err)
			continue
		}

		utils.Log("ok", "casino::data", "[UserTable] [Repair] added $1000 starting bonus to user:", wallet.UserID)
	}

	utils.Log("info", "casino::data", "[UserTable] [Repair] finished checking starting bonuses. Added to", len(wallets), "users.")
}
