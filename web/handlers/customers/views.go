package customers

import (
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/ctxfilter"
	"github/Doris-Mwito5/savannah-pos/internal/db"
	"github/Doris-Mwito5/savannah-pos/internal/dtos"
	"github/Doris-Mwito5/savannah-pos/internal/services"
	"github/Doris-Mwito5/savannah-pos/internal/utils"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func createCustomer(
	dB db.DB,
	customerService services.CustomerService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.CreateCustomerForm

		err := c.BindJSON(&req)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		log.Printf("DEBUG: Received data: %+v", req)
		log.Printf("DEBUG: Name: '%s', Email: '%s', CustomerType: '%v'", req.Name, req.Email, req.CustomerType)

		customer, err := customerService.CreateCustomer(c.Request.Context(), dB, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusCreated, customer)

	}
}

func updateCustomer(
	dB db.DB,
	customerService services.CustomerService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		var req dtos.UpdateCustomerForm

		err := c.BindJSON(&req)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		customer, err := customerService.UpdateCustomer(c.Request.Context(), dB, customerID, &req)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, customer)

	}
}

func getCustomer(
	dB db.DB,
	customerService services.CustomerService,
) func(c *gin.Context) {
	return func(c *gin.Context) {

		customerID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			appErr := apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
			utils.HandleError(c, appErr)
			return
		}

		customer, err := customerService.CustomerByID(c.Request.Context(), dB, customerID)
		if err != nil {
			utils.HandleError(c, err)
			return
		}

		c.JSON(http.StatusOK, customer)
	}
}

func listCustomers(
	dB db.DB,
	customerService services.CustomerService,
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

		customerList, err := customerService.ListShopCustomers(c.Request.Context(), dB, shopID, filter)
		if err != nil {
			utils.HandleError(c, err)
		}
		c.JSON(http.StatusOK, customerList)
	}
}
