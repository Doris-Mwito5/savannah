package products

import (
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/services"

	"github.com/gin-gonic/gin"
)

func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	productService services.ProductService,
) {
	r.POST("/products", createProduct(dB, productService))
	r.PUT("/products/:id", updateProduct(dB, productService))
	r.GET("/products/:id", getProduct(dB, productService))
	r.GET("/shop/:id/products", listProducts(dB, productService))
	r.DELETE("/products/:id", deleteProduct(dB, productService))
	r.GET("/categories/:id/average-price", getAveragePriceByCategory(dB, productService))
}
