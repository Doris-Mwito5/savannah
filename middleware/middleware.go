package middleware

import (
	"github/Doris-Mwito5/savannah-pos/auth"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware enforces OIDC token verification with better error handling
func AuthMiddleware(oidcSvc auth.OIDCService) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Printf("Auth middleware called for path: %s", c.FullPath())

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Printf("Missing Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "Missing Authorization header",
			})
			return
		}

		// Expect header in format: "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			log.Printf("Invalid Authorization header format")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "Invalid Authorization header format. Expected 'Bearer <token>'",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Printf("Empty token in Authorization header")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Authentication required",
				"message": "Empty token in Authorization header",
			})
			return
		}

		log.Printf("Verifying token...")

		userInfo, err := oidcSvc.VerifyToken(c.Request.Context(), token)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid or expired token",
				"message": "Token verification failed",
			})
			return
		}

		log.Printf("Token verified successfully for user: %s", userInfo.Email)

		// Store user info in context so handlers can use it
		c.Set("user", userInfo)
		c.Set("user_email", userInfo.Email)
		c.Set("user_sub", userInfo.Sub)
		c.Set("user_name", userInfo.Name)

		c.Next()
	}
}

// OptionalAuthMiddleware allows both authenticated and non-authenticated requests
func OptionalAuthMiddleware(oidcSvc auth.OIDCService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without authentication
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			// Invalid format, continue without authentication
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			// Empty token, continue without authentication
			c.Next()
			return
		}

		// Try to verify token, but don't fail if it doesn't work
		userInfo, err := oidcSvc.VerifyToken(c.Request.Context(), token)
		if err != nil {
			log.Printf("Optional auth failed: %v", err)
			// Continue without authentication
			c.Next()
			return
		}

		// Store user info if authentication succeeded
		c.Set("user", userInfo)
		c.Set("user_email", userInfo.Email)
		c.Set("user_sub", userInfo.Sub)
		c.Set("user_name", userInfo.Name)

		c.Next()
	}
}

// Helper function to get user from context
func GetUserFromContext(c *gin.Context) (*auth.OIDCUserInfo, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userInfo, ok := user.(*auth.OIDCUserInfo)
	return userInfo, ok
}

// Helper function for min
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}