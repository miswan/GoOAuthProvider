package models

import (
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex;not null"`
	UserID    uint   `gorm:"not null"`
	ClientID  string `gorm:"not null"`
	ExpiresAt time.Time
}
