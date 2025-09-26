package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github/Doris-Mwito5/savannah-pos/auth"
	"github/Doris-Mwito5/savannah-pos/config"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/domain"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/notification"
	"github/Doris-Mwito5/savannah-pos/internal/processor"
	"github/Doris-Mwito5/savannah-pos/internal/services"
	"github/Doris-Mwito5/savannah-pos/middleware"
	authhandler "github/Doris-Mwito5/savannah-pos/web/handlers/auth"
	"github/Doris-Mwito5/savannah-pos/web/handlers/categories"
	"github/Doris-Mwito5/savannah-pos/web/handlers/customers"
	"github/Doris-Mwito5/savannah-pos/web/handlers/orders"
	"github/Doris-Mwito5/savannah-pos/web/handlers/products"
)

type AppRouter struct {
	*gin.Engine
}

func BuildRouter(
	dB db.DB,
	domainStore *domain.Store,
) *AppRouter {
	router := gin.Default()

	baseRouter := router.Group("")
	baseAPIGroup := baseRouter.Group("/v1")

	baseAPIGroup.Use(middleware.CORSMiddleware())

	// Instantiate SMS notification services
	smsService := &models.SMSService{
		ApiKey:   config.AppConfig.SMSService.ApiKey,   
		Username: config.AppConfig.SMSService.Username, 
		BaseURL:  config.AppConfig.SMSService.BaseURL,
		Env:      config.AppConfig.SMSService.Env,
	}

	emailService := &models.EmailService{
        SMTPHost:   config.AppConfig.EmailService.SMTPHost,
        SMTPPort:   config.AppConfig.EmailService.SMTPPort,
        Username:   config.AppConfig.EmailService.Username,
        Password:   config.AppConfig.EmailService.Password,
        FromEmail:  config.AppConfig.EmailService.FromEmail,
        FromName:   config.AppConfig.EmailService.FromName,
        AdminEmail: config.AppConfig.EmailService.AdminEmail,
    }

	emailClient := processor.NewEmailClient(emailService)

	httpClient := &http.Client{} 
	smsClient := processor.NewSMSClient(smsService, httpClient)
	orderNotification := notification.NewOrderNotification(smsClient, emailClient)

	// Instantiate other services
	categoryService := services.NewCategoryService(domainStore)
	customerService := services.NewCustomerService(domainStore)
	// Pass the orderNotification to the OrderService
	orderService := services.NewOrderService(customerService, domainStore, orderNotification)
	productService := services.NewProductService(domainStore)

	// OIDC Auth service (now using config from .env)
	oidcService, err := auth.NewOIDCProvider(&config.AppConfig.OIDC)
	if err != nil {
		panic(err)
	}

	// AUTHENTICATION SETUP NOTE:
	// The OIDC authentication middleware is fully implemented and ready to use.
	// To protect API endpoints, simply uncomment the lines below:
	//
	// protected := baseAPIGroup.Group("/api")
	// protected.Use(middleware.AuthMiddleware(oidcService))
	//
	// Then register protected endpoints with the 'protected' group instead of 'baseAPIGroup'.
	// Currently endpoints are open for demonstration purposes.


	// Register endpoints
	authhandler.AddEndpoints(baseAPIGroup, dB, oidcService, customerService, config.AppConfig.JWTSecret)
	categories.AddEndpoints(baseAPIGroup, dB, categoryService)
	customers.AddEndpoints(baseAPIGroup, dB, customerService)
	orders.AddEndpoints(baseAPIGroup, dB, orderService)
	products.AddEndpoints(baseAPIGroup, dB, productService)

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error_message": "Endpoint not found"})
	})

	return &AppRouter{router}
}
