package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
}

type UserLogin struct {
	Username   string `json:"username" form:"username" validate:"required"`
	Password   string `json:"password" form:"password" validate:"required"`
	ContinueTo string `json:"continue_to" form:"continue_to"`
}

type UserRegister struct {
	Username string `json:"username" form:"username" validate:"required"`
	Password string `json:"password" form:"password" validate:"required"`
	Email    string `json:"email" form:"email" validate:"required,email"`
}