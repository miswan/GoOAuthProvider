package models

import (
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ClientID     string   `gorm:"column:client_id;uniqueIndex:idx_client_id;not null"`
	Secret       string   `gorm:"column:secret;not null"`
	RedirectURIs []string `gorm:"column:redirect_uris;type:text[];serializer:json"`
	GrantTypes   []string `gorm:"column:grant_types;type:text[];serializer:json"`
}

func (Client) TableName() string {
	return "clients"
}

type ClientRegistration struct {
	RedirectURIs []string `json:"redirect_uris" validate:"required,min=1"`
}

type AuthorizationRequest struct {
	ClientID            string `query:"client_id" validate:"required"`
	RedirectURI         string `query:"redirect_uri" validate:"required,url"`
	ResponseType        string `query:"response_type" validate:"required,oneof=code"`
	State               string `query:"state"`
	CodeChallenge      string `query:"code_challenge" validate:"required"`
	CodeChallengeMethod string `query:"code_challenge_method" validate:"required,oneof=S256 plain"`
}

type TokenRequest struct {
	GrantType    string `json:"grant_type" validate:"required,oneof=authorization_code refresh_token"`
	Code         string `json:"code"`
	RedirectURI  string `json:"redirect_uri"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CodeVerifier string `json:"code_verifier" validate:"required_if=GrantType authorization_code"`
	RefreshToken string `json:"refresh_token"`
}