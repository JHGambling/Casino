package tables

import (
	"errors"
	"fmt"
	"jhgambling/protocol"
	"jhgambling/protocol/models"
)

// WalletTable provides table operations for the UserModel
type WalletTable struct {
	protocol.BaseTable
}

// NewWalletTable creates a new user table
func NewWalletTable() *WalletTable {
	return &WalletTable{
		BaseTable: protocol.BaseTable{
			ID:    "wallets",
			Model: &models.WalletModel{},
		},
	}
}

// Create creates a new wallet
func (t *WalletTable) Create(data interface{}) error {
	wallet, ok := data.(*models.WalletModel)
	if !ok {
		return errors.New("invalid data type: expected *models.WalletModel")
	}

	t.PushRecordChange("create", nil, wallet)

	return t.DB.Create(wallet).Error
}

// FindByID finds a wallet by ID
func (t *WalletTable) FindByID(id interface{}) (interface{}, error) {
	var wallet models.WalletModel
	result := t.DB.First(&wallet, id)
	return &wallet, result.Error
}

// FindAll retrieves all wallets with pagination
func (t *WalletTable) FindAll(limit, offset int) ([]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}

	var wallets []models.WalletModel
	result := t.DB.Limit(limit).Offset(offset).Find(&wallets)
	if result.Error != nil {
		return nil, result.Error
	}

	// Convert to []interface{}
	items := make([]interface{}, len(wallets))
	for i, wallet := range wallets {
		walletCopy := wallet // Create a copy to avoid references to the same object
		items[i] = &walletCopy
	}

	return items, nil
}

// CreateAsUser implements user-based permission check for creating a user
func (t *WalletTable) CreateAsUser(user models.UserModel, data interface{}) error {
	return errors.New("you cant create a new wallet")
}

// FindByIDAsUser retrieves a wallet by ID
func (t *WalletTable) FindByIDAsUser(user models.UserModel, id interface{}) (interface{}, error) {
	foundWallet, err := t.FindByID(id)
	if err != nil {
		return nil, err
	}

	return foundWallet, nil
}

// FindAllAsUser retrieves all wallets
func (t *WalletTable) FindAllAsUser(user models.UserModel, limit, offset int) ([]interface{}, error) {
	wallets, err := t.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	return wallets, nil
}

// UpdateAsUser modifies a wallet with permission check
func (t *WalletTable) UpdateAsUser(user models.UserModel, id interface{}, data interface{}) error {
	// Users can update their own wallet
	idFloat, ok := id.(float64)
	if !ok {
		fmt.Println("well still no")
		return errors.New("invalid ID format: expected float64")
	}
	walletID := uint(idFloat)

	if user.Wallet.ID != walletID && !user.IsAdmin {
		fmt.Println("permission error")
		return errors.New("permission denied: you can only update your own wallet data")
	}

	walletData, ok := data.(*models.WalletModel)
	if !ok {
		return errors.New("invalid data type: expected *models.WalletModel")
	}

	// Find existing wallet
	var existingWallet models.WalletModel
	if err := t.DB.First(&existingWallet, walletID).Error; err != nil {
		return err
	}

	// Preserve starting bonus status if it's true
	if existingWallet.ReceivedStartingBonus {
		walletData.ReceivedStartingBonus = true
	}

	return t.Update(walletID, walletData)
}

// DeleteAsUser removes a wallet with permission check
func (t *WalletTable) DeleteAsUser(user models.UserModel, id interface{}) error {
	// Only admins can delete wallets
	if !user.IsAdmin {
		return errors.New("permission denied: only admins can delete wallets")
	}

	return t.Delete(id)
}

// Update updates a wallet and triggers notifications
func (t *WalletTable) Update(id interface{}, data interface{}) error {
	walletData, ok := data.(*models.WalletModel)
	if !ok {
		return errors.New("invalid data type: expected *models.WalletModel")
	}

	// Update the wallet
	err := t.DB.Model(&models.WalletModel{}).Where("id = ?", id).Updates(walletData).Error
	if err != nil {
		return err
	}

	// Fetch the updated wallet to push changes
	var updatedWallet models.WalletModel
	if err := t.DB.First(&updatedWallet, id).Error; err != nil {
		return err
	}

	// Push the change notification
	t.PushRecordChange("update", id, &updatedWallet)

	return nil
}
