package categories

import (
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/ctxfilter"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/services"
	"github/Doris-Mwito5/savannah-pos/internal/utils"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func createCategory(
	dB db.DB,
	categoryService services.CategoryService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.CreateCategoryForm

		err := c.BindJSON(&req)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		category, err := categoryService.CreateCategory(c.Request.Context(), dB, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, category)
	}
}

func getCategory(
	dB db.DB,
	categoryService services.CategoryService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		categoryID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}
		category, err := categoryService.CategoryByID(c.Request.Context(), dB, categoryID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}
		c.JSON(http.StatusOK, category)
	}

}

func listCategories(
	dB db.DB,
	categoryService services.CategoryService,
) func (c *gin.Context) {
	return func(c *gin.Context) {

		shopID := c.Param("id")
		if len(strings.TrimSpace(shopID)) < 1{
			appErr := apperr.NewBadRequest("shop id not set")
			utils.HandleError(c, appErr)
			return
		}

		filter, err := ctxfilter.FilterFromContext(c)
		if err != nil {
			appError := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appError)
			return
		}
		categoryList, err := categoryService.ListCategories(c.Request.Context(), dB, shopID, filter)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, categoryList)		
	}
}