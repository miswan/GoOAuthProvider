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

func (s *OAuthService) GenerateAuthorizationCode(clientID, redirectURI string, userID uint, codeChallenge, codeChallengeMethod string) (string, error) {
	code := utils.GenerateRandomString(32)
	err := s.store.StoreAuthCodeWithPKCE(code, clientID, redirectURI, userID, codeChallenge, codeChallengeMethod)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (s *OAuthService) ExchangeToken(req *models.TokenRequest) (string, string, error) {
	if req.GrantType != "authorization_code" && req.GrantType != "refresh_token" {
		return "", "", errors.New("unsupported grant type")
	}

	if req.GrantType == "authorization_code" {
		return s.handleAuthorizationCodeGrant(req)
	}

	return s.handleRefreshTokenGrant(req)
}

func (s *OAuthService) handleAuthorizationCodeGrant(req *models.TokenRequest) (string, string, error) {
	// Validate Client Credentials
	client := s.store.GetClient(req.ClientID)
	if client == nil {
		return "", "", errors.New("invalid client")
	}

	// If client has a secret, we must validate it.
	// In our model, all clients have secrets, so we enforce it.
	if client.Secret != "" && client.Secret != req.ClientSecret {
		return "", "", errors.New("invalid client secret")
	}

	authCode := s.store.GetAuthCode(req.Code)
	if authCode == nil {
		return "", "", errors.New("invalid authorization code")
	}

	if authCode.ClientID != req.ClientID {
		return "", "", errors.New("client id mismatch")
	}

	if authCode.RedirectURI != req.RedirectURI {
		return "", "", errors.New("redirect uri mismatch")
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

	// Validate Client Credentials for refresh token flow as well
	client := s.store.GetClient(req.ClientID)
	if client == nil {
		return "", "", errors.New("invalid client")
	}
	if client.Secret != "" && client.Secret != req.ClientSecret {
		return "", "", errors.New("invalid client secret")
	}

	if refreshToken.ClientID != req.ClientID {
		return "", "", errors.New("client id mismatch")
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

	// Remove padding if present in computedChallenge, although RawURLEncoding shouldn't have it.
	// But just in case standard encoding was used elsewhere.
	// Base64 RawURLEncoding is what we want.

	if !strings.EqualFold(computedChallenge, authCode.CodeChallenge) {
		return errors.New("invalid code verifier")
	}

	return nil
}
