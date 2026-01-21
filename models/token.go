package models

import (
	"gorm.io/gorm"
	"time"
)

type AuthCode struct {
	gorm.Model
	Code               string `gorm:"uniqueIndex;not null"`
	ClientID           string `gorm:"not null"`
	UserID             uint   `gorm:"not null"`
	RedirectURI        string `gorm:"not null"`
	ExpiresAt          time.Time
	CodeChallenge      string
	CodeChallengeMethod string
	Used               bool
}

type RefreshToken struct {
	gorm.Model
	Token    string `gorm:"uniqueIndex;not null"`
	UserID   uint   `gorm:"not null"`
	ClientID string `gorm:"not null"`
	ExpiresAt time.Time
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
	GrantType    string `json:"grant_type" form:"grant_type" validate:"required,oneof=authorization_code refresh_token"`
	Code         string `json:"code" form:"code"`
	RedirectURI  string `json:"redirect_uri" form:"redirect_uri"`
	ClientID     string `json:"client_id" form:"client_id"`
	ClientSecret string `json:"client_secret" form:"client_secret"`
	CodeVerifier string `json:"code_verifier" form:"code_verifier" validate:"required_if=GrantType authorization_code"`
	RefreshToken string `json:"refresh_token" form:"refresh_token"`
}
