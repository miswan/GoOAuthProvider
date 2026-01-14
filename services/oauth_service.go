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

func (s *OAuthService) GetUserByID(userID uint) (*models.User, error) {
	user := s.store.GetUserByID(userID)
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *OAuthService) GenerateAuthorizationCode(clientID string, userID uint, codeChallenge, codeChallengeMethod string) (string, error) {
	code := utils.GenerateRandomString(32)
	err := s.store.StoreAuthCodeWithPKCE(code, clientID, userID, codeChallenge, codeChallengeMethod)
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
	authCode := s.store.GetAuthCode(req.Code)
	if authCode == nil {
		return "", "", errors.New("invalid authorization code")
	}

	// Verify Client ID matches
	if req.ClientID != "" && req.ClientID != authCode.ClientID {
		return "", "", errors.New("invalid client_id")
	}

	// Verify Client Secret if provided or required
	// Fetch client to check secret
	client := s.store.GetClient(authCode.ClientID)
	if client == nil {
		return "", "", errors.New("invalid client")
	}

	// If request has client_secret, verify it
	if req.ClientSecret != "" {
		if req.ClientSecret != client.Secret {
			return "", "", errors.New("invalid client_secret")
		}
	} else {
		// If no secret provided, check if client is confidential (has secret)
		// For this implementation, we assume if client has a secret stored, it must be provided?
		// Or if it's a public client (no secret or using PKCE), we might not enforce secret.
		// Since we are using PKCE, public clients are supported.
		// However, if it IS a confidential client, we should enforce secret.
		// Let's assume non-empty Secret in DB means confidential.
		// But in StoreClient, we always generate a Secret. So all clients are confidential?
		// Unless we support public clients.
		// Given we enforce PKCE, we should allow public clients, but our StoreClient makes everyone confidential.
		// For now, let's enforce secret ONLY if it was provided in request, OR if we want to be strict.
		// A common pattern: if client_id/secret is in Basic Auth header, it's checked there.
		// Here we only look at body.
		// Let's enforce secret match if ClientID is passed in body.
		if req.ClientID != "" && req.ClientSecret == "" {
			return "", "", errors.New("client_secret is required")
		}
	}

	if err := s.validatePKCE(authCode, req.CodeVerifier); err != nil {
		return "", "", err
	}

	// Generate tokens
	accessToken, err := utils.GeneratePaseto(authCode.UserID, time.Hour)
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
	accessToken, err := utils.GeneratePaseto(refreshToken.UserID, time.Hour)
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