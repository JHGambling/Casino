package models

import (
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model

	Username     string
	DisplayName  string
	PasswordHash string
	JoinedAt     time.Time
	IsAdmin      bool
}
