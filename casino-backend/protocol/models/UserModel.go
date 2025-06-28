package models

import (
	"time"

	"gorm.io/gorm"
)

// UserModel represents a user in the system.
type UserModel struct {
	gorm.Model

	Username     string
	DisplayName  string
	PasswordHash string
	JoinedAt     time.Time
	IsAdmin      bool

	Wallet WalletModel `gorm:"foreignKey:UserID"`
}
