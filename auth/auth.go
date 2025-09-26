package auth

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github/Doris-Mwito5/savannah-pos/config"
)

// OIDCProvider implements OIDCService
type OIDCProvider struct {
	provider     *oidc.Provider
	oauth2Config oauth2.Config
	verifier     *oidc.IDTokenVerifier
	clientID     string
	clientSecret string
	redirectURL  string
	issuerURL    string
}

type OIDCService interface {
	GetAuthURL(state string) string
	HandleCallback(ctx context.Context, code, state string) (*OIDCUserInfo, error)
	VerifyToken(ctx context.Context, rawIDToken string) (*OIDCUserInfo, error)
	CreateAuthURL(c *gin.Context) (string, string, error)
}

type OIDCUserInfo struct {
	Sub           string `json:"sub"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
	PhoneNumber   string `json:"phone_number"`
}

func NewOIDCProvider(cfg *config.OIDCConfig) (OIDCService, error) {
	if err := validateOIDCConfig(cfg); err != nil {
		return nil, fmt.Errorf("invalid OIDC configuration: %w", err)
	}

	ctx := context.Background()

	log.Printf("Initializing OIDC provider with issuer: %s", cfg.IssuerURL)

	ctxWithTimeout, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	provider, err := oidc.NewProvider(ctxWithTimeout, cfg.IssuerURL)
	if err != nil {
		log.Printf("Failed to create OIDC provider: %v", err)
		return nil, fmt.Errorf("failed to create OIDC provider for %s: %w", cfg.IssuerURL, err)
	}

	log.Printf("OIDC provider created successfully")

	oauth2Config := oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		RedirectURL:  cfg.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"}, 
	}

	verifier := provider.Verifier(&oidc.Config{
		ClientID: cfg.ClientID,
	})

	log.Printf("OAuth2 config: ClientID=%s, RedirectURL=%s", cfg.ClientID, cfg.RedirectURL)

	return &OIDCProvider{
		provider:     provider,
		oauth2Config: oauth2Config,
		verifier:     verifier,
		clientID:     cfg.ClientID,
		clientSecret: cfg.ClientSecret,
		redirectURL:  cfg.RedirectURL,
		issuerURL:    cfg.IssuerURL,
	}, nil
}

func validateOIDCConfig(cfg *config.OIDCConfig) error {
	if cfg.ClientID == "" {
		return errors.New("client ID is required")
	}
	if cfg.ClientSecret == "" {
		return errors.New("client secret is required")
	}
	if cfg.IssuerURL == "" {
		return errors.New("issuer URL is required")
	}
	if cfg.RedirectURL == "" {
		return errors.New("redirect URL is required")
	}
	return nil
}

func (o *OIDCProvider) CreateAuthURL(c *gin.Context) (string, string, error) {
	state, err := generateRandomState()
	if err != nil {
		log.Printf("Failed to generate state: %v", err)
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}

	log.Printf("Generated state: %s", state)

	c.SetCookie(
		"oauth_state",    // name
		state,           // value
		600,            // maxAge (10 minutes)
		"/",            // path
		"",             // domain (empty for current domain)
		false,          // secure (set to true in production with HTTPS)
		true,           // httpOnly
	)

	authURL := o.GetAuthURL(state)
	log.Printf("Generated auth URL: %s", authURL)

	return authURL, state, nil
}

// GetAuthURL returns the OIDC login URL
func (o *OIDCProvider) GetAuthURL(state string) string {
	return o.oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// HandleCallback exchanges the code for tokens and verifies ID token
func (o *OIDCProvider) HandleCallback(ctx context.Context, code, state string) (*OIDCUserInfo, error) {
	log.Printf("Handling callback with code: %s, state: %s", code[:min(len(code), 10)]+"...", state)

	// Exchange authorization code for tokens
	token, err := o.oauth2Config.Exchange(ctx, code)
	if err != nil {
		log.Printf("Failed to exchange token: %v", err)
		return nil, fmt.Errorf("failed to exchange token: %w", err)
	}

	log.Printf("Token exchange successful")

	// Extract ID token from OAuth2 token
	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		log.Printf("No id_token in OAuth2 response")
		return nil, errors.New("no id_token field in oauth2 token")
	}

	log.Printf("ID token found, verifying...")

	// Verify and parse the ID token
	return o.VerifyToken(ctx, rawIDToken)
}

// VerifyToken validates the ID token and extracts claims
func (o *OIDCProvider) VerifyToken(ctx context.Context, rawIDToken string) (*OIDCUserInfo, error) {
	log.Printf("Verifying ID token...")

	idToken, err := o.verifier.Verify(ctx, rawIDToken)
	if err != nil {
		log.Printf("Failed to verify ID token: %v", err)
		return nil, fmt.Errorf("failed to verify ID token: %w", err)
	}

	log.Printf("ID token verified successfully")

	var userInfo OIDCUserInfo
	if err := idToken.Claims(&userInfo); err != nil {
		log.Printf("Failed to parse ID token claims: %v", err)
		return nil, fmt.Errorf("failed to parse ID token claims: %w", err)
	}

	log.Printf("User info extracted: email=%s, name=%s, sub=%s", userInfo.Email, userInfo.Name, userInfo.Sub)

	return &userInfo, nil
}

// generateRandomState builds a secure OAuth2 state string
func generateRandomState() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}