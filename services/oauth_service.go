package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
	"strings"
	"time"
)

type OAuthService struct {
	store *storage.MemoryStorage
}

func NewOAuthService(store *storage.MemoryStorage) *OAuthService {
	return &OAuthService{store: store}
}

func (s *OAuthService) ValidateAuthorizationRequest(req *models.AuthorizationRequest) error {
	client := s.store.GetClient(req.ClientID)
	if client == nil {
		return errors.New("invalid client")
	}

	// Validate redirect URI
	validURI := false
	for _, uri := range client.RedirectURIs {
		if uri == req.RedirectURI {
			validURI = true
			break
		}
	}
	if !validURI {
		return errors.New("invalid redirect URI")
	}

	// Validate PKCE parameters
	if req.CodeChallenge == "" {
		return errors.New("code_challenge is required")
	}
	if req.CodeChallengeMethod != "S256" && req.CodeChallengeMethod != "plain" {
		return errors.New("code_challenge_method must be 'S256' or 'plain'")
	}

	return nil
}

func (s *OAuthService) GenerateAuthorizationCode(clientID, userID string, codeChallenge, codeChallengeMethod string) (string, error) {
	code := utils.GenerateRandomString(32)
	s.store.StoreAuthCodeWithPKCE(code, clientID, userID, codeChallenge, codeChallengeMethod)
	return code, nil
}

func (s *OAuthService) ExchangeToken(req *models.TokenRequest) (string, string, error) {
	// Validate grant type
	if req.GrantType != "authorization_code" && req.GrantType != "refresh_token" {
		return "", "", errors.New("unsupported grant type")
	}

	if req.GrantType == "authorization_code" {
		return s.handleAuthorizationCodeGrant(req)
	}

	return s.handleRefreshTokenGrant(req)
}

func (s *OAuthService) handleAuthorizationCodeGrant(req *models.TokenRequest) (string, string, error) {
	// Validate authorization code
	authCode := s.store.GetAuthCode(req.Code)
	if authCode == nil {
		return "", "", errors.New("invalid authorization code")
	}

	// Validate code verifier with PKCE
	if err := s.validatePKCE(authCode, req.CodeVerifier); err != nil {
		return "", "", err
	}

	// Generate tokens
	accessToken, err := utils.GenerateJWT(authCode.UserID, time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken := utils.GenerateRandomString(32)
	s.store.StoreRefreshToken(refreshToken, authCode.UserID)

	return accessToken, refreshToken, nil
}

// Implementasi handleRefreshTokenGrant
func (s *OAuthService) handleRefreshTokenGrant(req *models.TokenRequest) (string, string, error) {
	if req.RefreshToken == "" {
		return "", "", errors.New("refresh token is required")
	}

	userID := s.store.GetRefreshToken(req.RefreshToken)
	if userID == "" {
		return "", "", errors.New("invalid refresh token")
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(userID, time.Hour)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshToken := utils.GenerateRandomString(32)
	s.store.StoreRefreshToken(newRefreshToken, userID)

	return accessToken, newRefreshToken, nil
}

func (s *OAuthService) validatePKCE(authCode *storage.AuthCode, codeVerifier string) error {
	if authCode.CodeChallenge == "" {
		return errors.New("code challenge not found")
	}

	var computedChallenge string
	if authCode.CodeChallengeMethod == "S256" {
		h := sha256.New()
		h.Write([]byte(codeVerifier))
		computedChallenge = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	} else { // plain
		computedChallenge = codeVerifier
	}

	if !strings.EqualFold(computedChallenge, authCode.CodeChallenge) {
		return errors.New("invalid code verifier")
	}

	return nil
}