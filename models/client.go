package models

import (
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ClientID     string         `gorm:"uniqueIndex;not null"`
	Secret       string         `gorm:"not null"`
	RedirectURIs pq.StringArray `gorm:"type:text[]"`
	GrantTypes   pq.StringArray `gorm:"type:text[]"`
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