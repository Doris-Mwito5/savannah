package middleware

import (
	"github/Doris-Mwito5/savannah-pos/auth"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

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

func OptionalAuthMiddleware(oidcSvc auth.OIDCService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.Next()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			c.Next()
			return
		}

		userInfo, err := oidcSvc.VerifyToken(c.Request.Context(), token)
		if err != nil {
			log.Printf("Optional auth failed: %v", err)
			c.Next()
			return
		}

		c.Set("user", userInfo)
		c.Set("user_email", userInfo.Email)
		c.Set("user_sub", userInfo.Sub)
		c.Set("user_name", userInfo.Name)

		c.Next()
	}
}

func GetUserFromContext(c *gin.Context) (*auth.OIDCUserInfo, bool) {
	user, exists := c.Get("user")
	if !exists {
		return nil, false
	}

	userInfo, ok := user.(*auth.OIDCUserInfo)
	return userInfo, ok
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}