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
	store *storage.PostgresStorage
}

func NewOAuthService(store *storage.PostgresStorage) *OAuthService {
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

func (s *OAuthService) GenerateAuthorizationCode(clientID string, userID uint, codeChallenge, codeChallengeMethod, redirectURI string) (string, error) {
	code := utils.GenerateRandomString(32)
	err := s.store.StoreAuthCodeWithPKCE(code, clientID, redirectURI, userID, codeChallenge, codeChallengeMethod)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *OAuthService) ExchangeToken(req *models.TokenRequest) (string, string, error) {
	// Authenticate client
	// If the client exists and has a secret, we MUST validate it (Confidential Client).
	// If the client exists and has NO secret, it's a Public Client (we don't check secret).
	// But in this simple implementation, we assume all clients in DB have secrets.
	// However, we should check if the client exists first.
	client := s.store.GetClient(req.ClientID)
	if client == nil {
		return "", "", errors.New("invalid client")
	}

	// If client has a secret stored, we must validate it.
	// If the request provides a secret, it must match.
	// If the request does NOT provide a secret, we must check if the client requires one.
	if client.Secret != "" {
		if req.ClientSecret == "" || req.ClientSecret != client.Secret {
			return "", "", errors.New("invalid client credentials")
		}
	}

	if req.GrantType != "authorization_code" && req.GrantType != "refresh_token" {
		return "", "", errors.New("unsupported grant type")
	}

	if req.GrantType == "authorization_code" {
		return s.handleAuthorizationCodeGrant(req)
	}

	return s.handleRefreshTokenGrant(req)
}

func (s *OAuthService) handleAuthorizationCodeGrant(req *models.TokenRequest) (string, string, error) {
	authCode := s.store.GetAuthCode(req.Code)
	if authCode == nil {
		return "", "", errors.New("invalid authorization code")
	}

	// Validate ClientID
	if authCode.ClientID != req.ClientID {
		return "", "", errors.New("invalid client id")
	}

	// Validate RedirectURI
	if authCode.RedirectURI != req.RedirectURI {
		return "", "", errors.New("invalid redirect uri")
	}

	if err := s.validatePKCE(authCode, req.CodeVerifier); err != nil {
		return "", "", err
	}

	// Generate tokens
	accessToken, err := utils.GenerateJWT(authCode.UserID, time.Hour)
	if err != nil {
		return "", "", err
	}

	refreshToken := utils.GenerateRandomString(32)
	err = s.store.StoreRefreshToken(refreshToken, authCode.UserID, authCode.ClientID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (s *OAuthService) handleRefreshTokenGrant(req *models.TokenRequest) (string, string, error) {
	if req.RefreshToken == "" {
		return "", "", errors.New("refresh token is required")
	}

	refreshToken := s.store.GetRefreshToken(req.RefreshToken)
	if refreshToken == nil {
		return "", "", errors.New("invalid refresh token")
	}

	// Delete the used refresh token
	if err := s.store.DeleteRefreshToken(req.RefreshToken); err != nil {
		return "", "", err
	}

	// Generate new access token
	accessToken, err := utils.GenerateJWT(refreshToken.UserID, time.Hour)
	if err != nil {
		return "", "", err
	}

	// Generate new refresh token
	newRefreshToken := utils.GenerateRandomString(32)
	err = s.store.StoreRefreshToken(newRefreshToken, refreshToken.UserID, refreshToken.ClientID)
	if err != nil {
		return "", "", err
	}

	return accessToken, newRefreshToken, nil
}

func (s *OAuthService) validatePKCE(authCode *models.AuthCode, codeVerifier string) error {
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