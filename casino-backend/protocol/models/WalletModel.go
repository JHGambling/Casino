package models

import (
	"gorm.io/gorm"
)

// The WalletModel represents the wallet of one player
type WalletModel struct {
	gorm.Model

	UserID uint

	ReceivedStartingBonus bool
	NetworthCents         uint
}
