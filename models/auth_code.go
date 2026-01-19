package models

import (
	"gorm.io/gorm"
	"time"
)

type AuthCode struct {
	gorm.Model
	Code                string `gorm:"uniqueIndex;not null"`
	ClientID            string `gorm:"not null"`
	UserID              uint   `gorm:"not null"`
	RedirectURI         string `gorm:"not null"` // Added as per memory
	ExpiresAt           time.Time
	CodeChallenge       string
	CodeChallengeMethod string
	Used                bool `gorm:"default:false"`
}

type RefreshToken struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex;not null"`
	UserID    uint   `gorm:"not null"`
	ClientID  string `gorm:"not null"`
	ExpiresAt time.Time
}
