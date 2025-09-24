package categories

import (
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/services"

	"github.com/gin-gonic/gin"
)


func AddEndpoints(
	r *gin.RouterGroup,
	dB db.DB,
	categoryService services.CategoryService,
) {
	r.POST("/categories", createCategory(dB, categoryService))
	r.GET("/categories/:id", getCategory(dB, categoryService))
	r.GET("shop/:id/categories", listCategories(dB, categoryService))
}
