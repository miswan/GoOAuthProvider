package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ClientID     string         `gorm:"column:client_id;uniqueIndex:idx_client_id;not null"`
	Secret       string         `gorm:"column:secret;not null"`
	RedirectURIs pq.StringArray `gorm:"column:redirect_uris;type:text[]"`
	GrantTypes   pq.StringArray `gorm:"column:grant_types;type:text[]"`
}

func (Client) TableName() string {
	return "clients"
}

type ClientRegistration struct {
	RedirectURIs []string `json:"redirect_uris" validate:"required,min=1"`
}
