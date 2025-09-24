package products

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

func createProduct(
	dB db.DB,
	productService services.ProductService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.CreateProductForm

		err := c.BindJSON(&req)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		product, err := productService.CreateProduct(c.Request.Context(), dB, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, product)

	}
}

func updateProduct(
	dB db.DB,
	productService services.ProductService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.UpdateProductForm

		err := c.BindJSON(&req)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		product, err := productService.UpdateProduct(c.Request.Context(), dB, productID, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, product)

	}
}

func getProduct(
	dB db.DB,
	productService services.ProductService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		product, err := productService.ProductByID(c.Request.Context(), dB, productID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func listProducts(
	dB db.DB,
	productService services.ProductService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		shopID := c.Param("id")
		if len(strings.TrimSpace(shopID)) < 1 {
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

		productList, err := productService.ListProducts(c.Request.Context(), dB, shopID, filter)
		if err != nil {
			utils.HandleError(c, err)
		}
		c.JSON(http.StatusOK, productList)
	}
}

func deleteProduct(
	dB db.DB,
	productService services.ProductService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		productID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		product, err := productService.DeleteProduct(c.Request.Context(), dB, productID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, product)
	}
}

func getAveragePriceByCategory(
	dB db.DB,
	productService services.ProductService,
) func(c *gin.Context) {
	return func(c *gin.Context) {
		categoryID, err := strconv.ParseInt(c.Param("category_id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		averagePrice, err := productService.GetAveragePriceByCategory(c.Request.Context(), dB, categoryID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, averagePrice)
	}
}
