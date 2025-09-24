package auth

import (
	"github/Doris-Mwito5/savannah-pos/auth"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/services"

	"github.com/gin-gonic/gin"
)

// AddEndpoints registers authentication-related routes
func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	oidcService auth.OIDCService,
	customerService services.CustomerService,
	jwtSecret string,
) {
	r.GET("/auth/login", login(dB, oidcService))
	r.GET("/auth/callback", callback(dB, oidcService, customerService, jwtSecret))
	r.POST("/auth/logout", logout(dB))

	// Protected routes
	authGroup := r.Group("/auth")
	authGroup.Use() 
}
