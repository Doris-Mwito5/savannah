package customers

import (
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/services"

	"github.com/gin-gonic/gin"
)

func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	customerService services.CustomerService,
) {
	r.POST("/customers", createCustomer(dB, customerService))
	r.PUT("/customers/:id", updateCustomer(dB, customerService))
	r.GET("/customers/:id", getCustomer(dB, customerService))
	r.GET("/shop/:id/customers", listCustomers(dB, customerService))
}
