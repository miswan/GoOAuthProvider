package models

import (
	"gorm.io/gorm"
	"time"
)

type AuthCode struct {
	gorm.Model
	Code                string
	ClientID            string
	UserID              uint
	RedirectURI         string // Added RedirectURI
	Scope               string // Added Scope
	ExpiresAt           time.Time
	CodeChallenge       string
	CodeChallengeMethod string
	Used                bool
}

type RefreshToken struct {
	gorm.Model
	Token     string `gorm:"uniqueIndex;not null"`
	UserID    uint   `gorm:"not null"`
	ClientID  string `gorm:"not null"`
	Scope     string // Added Scope
	ExpiresAt time.Time
}
