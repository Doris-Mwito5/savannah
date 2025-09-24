package orders

import (
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/services"

	"github.com/gin-gonic/gin"
)

func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	orderService services.OrderService,
) {
	r.POST("/orders", createOrder(dB, orderService))
	// r.PUT("/orders/:id", updateOrder(dB, orderService))
	r.GET("/orders/:id", getOrder(dB, orderService))
	r.GET("/shop/:id/orders", listOrders(dB, orderService))

}
