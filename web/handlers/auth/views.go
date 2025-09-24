package auth

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github/Doris-Mwito5/savannah-pos/auth"
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/services"
	"github/Doris-Mwito5/savannah-pos/internal/utils"
)

// JWT Claims structure
type Claims struct {
	CustomerID int64  `json:"customer_id"`
	Email      string `json:"email"`
	Sub        string `json:"sub"`
	jwt.RegisteredClaims
}

func login(
	dB db.DB,
	oidcService auth.OIDCService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Printf("Login request received from IP: %s", c.ClientIP())

		authURL, state, err := oidcService.CreateAuthURL(c)
		if err != nil {
			log.Printf("Failed to create auth URL: %v", err)
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		log.Printf("Generated auth URL for login")

		// Check if client wants JSON response or browser redirect
		if c.GetHeader("Accept") == "application/json" || c.Query("format") == "json" {
			c.JSON(http.StatusOK, gin.H{
				"auth_url": authURL,
				"state":    state,
				"message":  "Redirect to auth_url to authenticate",
			})
			return
		}

		// Direct redirect for browser clients
		c.Redirect(http.StatusFound, authURL)
	}
}

func callback(
	dB db.DB,
	oidcService auth.OIDCService,
	customerService services.CustomerService,
	jwtSecret string,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		code := c.Query("code")
		state := c.Query("state")
		errorParam := c.Query("error")
		errorDesc := c.Query("error_description")

		log.Printf("Callback received - Code present: %t, State: %s, Error: %s",
			code != "", state, errorParam)

		// Handle OAuth errors
		if errorParam != "" {
			log.Printf("OAuth error: %s - %s", errorParam, errorDesc)
			appErr := apperr.NewBadRequest("OAuth authentication failed: " + errorParam)
			utils.HandleError(c, appErr)
			return
		}

		if code == "" || state == "" {
			appErr := apperr.NewBadRequest("missing authorization code or state")
			utils.HandleError(c, appErr)
			return
		}

		// Verify state parameter from cookie
		savedState, err := c.Cookie("oauth_state")
		if err != nil || savedState != state {
			log.Printf("State mismatch. Expected: %s, Got: %s", savedState, state)
			appErr := apperr.NewBadRequest("invalid or missing state parameter")
			utils.HandleError(c, appErr)
			return
		}

		// Clear state cookie
		c.SetCookie("oauth_state", "", -1, "/", "", false, true)

		// Exchange code for tokens and user info
		userInfo, err := oidcService.HandleCallback(c.Request.Context(), code, state)
		if err != nil {
			log.Printf("Failed to handle callback: %v", err)
			utils.HandleError(c, apperr.NewBadRequest("failed to complete authentication"))
			return
		}

		log.Printf("Authentication successful for user: %s", userInfo.Email)

		// Try to fetch customer by email
		customer, err := customerService.CustomerByEmail(c.Request.Context(), dB, userInfo.Email)
		if err != nil {
			log.Printf("Customer not found, creating new: %s", userInfo.Email)

			newCustomer := &dtos.CreateCustomerForm{
				Email: userInfo.Email,
				Name:  userInfo.Name,
			}

			customer, err = customerService.CreateCustomer(c.Request.Context(), dB, newCustomer)
			if err != nil {
				log.Printf("Failed to create customer: %v", err)
				utils.HandleError(c, apperr.NewInternal("could not create customer"))
				return
			}
		}

		// Generate JWT token
		token, expiresAt, err := generateJWT(customer, userInfo, jwtSecret)
		if err != nil {
			log.Printf("Failed to generate JWT: %v", err)
			utils.HandleError(c, apperr.NewInternal("failed to generate token"))
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message":    "Authentication successful",
			"token":      token,
			"expires_at": expiresAt,
			"customer":   customer,
			"user_info":  userInfo,
		})
	}
}

func generateJWT(customer *models.Customer, userInfo *auth.OIDCUserInfo, jwtSecret string) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := Claims{
		CustomerID: customer.ID,
		Email:      userInfo.Email,
		Sub:        userInfo.Sub,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Subject:   userInfo.Sub,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(jwtSecret))
	return signed, expiresAt, err
}

func logout(
	dB db.DB,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		log.Printf("Logout request from IP: %s", c.ClientIP())

		// In production: invalidate JWT or clear session

		c.JSON(http.StatusOK, gin.H{
			"message": "Logged out successfully",
		})
	}
}
