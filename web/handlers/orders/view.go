package orders

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

func createOrder(
	dB db.DB,
	orderService services.OrderService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.CreateOrderForm

		err := c.BindJSON(&req)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		order, err := orderService.CreateOrder(c.Request.Context(), dB, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, order)

	}
}

// func updateOrder(
// 	dB db.DB,
// 	orderService services.OrderService,
// ) func(c *gin.Context) {
// 	return func(c *gin.Context) {

// 		var req dtos.UpdateOrderForm

// 		err := c.BindJSON(&req)
// 		if err != nil {
// 			appErr := apperr.NewErrorWithType(
// 				err,
// 				apperr.BadRequest,
// 			)
// 			utils.HandleError(c, appErr)
// 			return
// 		}

// 		orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
// 		if err != nil {
// 			appErr := apperr.NewErrorWithType(
// 				err,
// 				apperr.BadRequest,
// 			)
// 			utils.HandleError(c, appErr)
// 			return
// 		}

// 		order, err := orderService.UpdateProduct(c.Request.Context(), dB, orderID, &req)
// 		if err != nil {
// 			utils.HandleError(c, err)
// 			return
// 		}

// 		c.JSON(http.StatusOK, order)

// 	}
// }

func getOrder(
	dB db.DB,
	orderService services.OrderService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		orderID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		order, err := orderService.OrderByID(c.Request.Context(), dB, orderID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, order)
	}
}

func listOrders(
	dB db.DB,
	orderService services.OrderService,
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

		orderList, err := orderService.ListShopOrders(c.Request.Context(), dB, shopID, filter)
		if err != nil {
			utils.HandleError(c, err)
		}
		c.JSON(http.StatusOK, orderList)
	}
}
