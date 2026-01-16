package services

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"oauth2-provider/models"
	"oauth2-provider/storage"
	"oauth2-provider/utils"
	"time"
)

type OAuthService struct {
	store storage.Storage
}

func NewOAuthService(store storage.Storage) *OAuthService {
	return &OAuthService{store: store}
}

func (s *OAuthService) ValidateAuthorizationRequest(req *models.AuthorizationRequest) error {
	client := s.store.GetClient(req.ClientID)
	if client == nil {
		return errors.New("invalid client")
	}

	// Validate redirect URI
	// RFC 6749 Section 3.1.2.2: The authorization server MUST validate that the
	// redirect_uri provided matches a registered redirect URI.
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

func (s *OAuthService) GenerateAuthorizationCode(clientID string, userID uint, redirectURI, codeChallenge, codeChallengeMethod string) (string, error) {
	code := utils.GenerateRandomString(32)
	err := s.store.StoreAuthCodeWithPKCE(code, clientID, userID, redirectURI, codeChallenge, codeChallengeMethod)
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

	// Validate Redirect URI matches the one used in Authorization Request
	if authCode.RedirectURI != req.RedirectURI {
		return "", "", errors.New("redirect_uri mismatch")
	}

	// Verify client credentials if provided (Confidential Clients)
	// Or ensure the clientID matches the one in the auth code (Public Clients)
	if req.ClientID != "" && authCode.ClientID != req.ClientID {
		return "", "", errors.New("client_id mismatch")
	}

	if err := s.validatePKCE(authCode, req.CodeVerifier); err != nil {
		return "", "", err
	}

	// Generate tokens using Paseto
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

	// Delete the used refresh token (Rotation)
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
		// If code challenge was not stored, PKCE was not used.
		// If the client is public, this is bad, but for legacy it might pass.
		// However, we enforce PKCE in ValidateAuthorizationRequest, so this shouldn't happen.
		return errors.New("code challenge not found")
	}

	if codeVerifier == "" {
		return errors.New("code_verifier is required")
	}

	var computedChallenge string
	if authCode.CodeChallengeMethod == "S256" {
		h := sha256.New()
		h.Write([]byte(codeVerifier))
		computedChallenge = base64.RawURLEncoding.EncodeToString(h.Sum(nil))
	} else { // plain
		computedChallenge = codeVerifier
	}

	if computedChallenge != authCode.CodeChallenge {
		return errors.New("invalid code verifier")
	}

	return nil
}
