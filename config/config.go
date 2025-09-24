package config

import (
	"fmt"
	"github/Doris-Mwito5/savannah-pos/env"
)

// Config holds all runtime configuration
type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	SMSService  SMSServiceConfig // ✅ Better structure for SMS config
	EmailService EmailServiceConfig
	OIDC        OIDCConfig
}

// SMSServiceConfig holds Africa's Talking settings
type SMSServiceConfig struct {
	BaseURL  string
	ApiKey   string
	Username string
	Env      string
}

type EmailServiceConfig struct {
    SMTPHost   string
    SMTPPort   int
    Username   string
    Password   string
    FromEmail  string
    FromName   string
    AdminEmail string
}

// OIDCConfig holds OpenID Connect settings
type OIDCConfig struct {
	ClientID     string
	ClientSecret string
	IssuerURL    string
	RedirectURL  string
}

var AppConfig Config

// LoadEnvConfig reads configuration from env vars
func LoadEnvConfig() error {
	databaseURL, _ := env.GetEnvString("DATABASE_URL")
	port, _ := env.GetEnvString("PORT")
	jwtSecret, _ := env.GetEnvString("JWT_SECRET")

	// 🔑 Load Africa's Talking settings from environment variables
	baseURL, _ := env.GetEnvString("AFRICAS_TALKING_BASE_URL")
	apiKey, _ := env.GetEnvString("AFRICAS_TALKING_API_KEY")
	username, _ := env.GetEnvString("AFRICAS_TALKING_USERNAME")
	envVar, _ := env.GetEnvString("AFRICAS_TALKING_ENV") // ✅ FIXED: Changed variable name

	smtpHost, _ := env.GetEnvString("SMTP_HOST")
    smtpPort, _ := env.GetEnvInt("SMTP_PORT")
    smtpUsername, _ := env.GetEnvString("SMTP_USERNAME")
    smtpPassword, _ := env.GetEnvString("SMTP_PASSWORD")
    fromEmail, _ := env.GetEnvString("SMTP_FROM_EMAIL")
    fromName, _ := env.GetEnvString("SMTP_FROM_NAME")
    adminEmail, _ := env.GetEnvString("ADMIN_EMAIL")
	
	// OIDC settings
	clientID, _ := env.GetEnvString("OIDC_CLIENT_ID")
	clientSecret, _ := env.GetEnvString("OIDC_CLIENT_SECRET")
	issuerURL, _ := env.GetEnvString("OIDC_ISSUER_URL")
	redirectURL, _ := env.GetEnvString("OIDC_REDIRECT_URL")

	AppConfig = Config{
		DatabaseURL: databaseURL,
		Port:        port,
		JWTSecret:   jwtSecret,
		SMSService: SMSServiceConfig{
			BaseURL:  baseURL,
			ApiKey:   apiKey,
			Username: username,
			Env:      envVar, // ✅ FIXED: Use the renamed variable
		},
		EmailService: EmailServiceConfig{
            SMTPHost:   smtpHost,
            SMTPPort:   smtpPort,
            Username:   smtpUsername,
            Password:   smtpPassword,
            FromEmail:  fromEmail,
            FromName:   fromName,
            AdminEmail: adminEmail,
        },
		OIDC: OIDCConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			IssuerURL:    issuerURL,
			RedirectURL:  redirectURL,
		},
	}

	// Validate required SMS config
	if AppConfig.SMSService.ApiKey == "" || AppConfig.SMSService.Username == "" {
		return fmt.Errorf("Africa's Talking API credentials are required")
	}

	// Log the loaded config for debugging
	fmt.Printf("SMS Config Loaded - Username: %s, API Key: %s, BaseURL: %s, Env: %s\n", 
		AppConfig.SMSService.Username, 
		AppConfig.SMSService.ApiKey, 
		AppConfig.SMSService.BaseURL,
		AppConfig.SMSService.Env)

	return nil
}