package utils

import (
	"errors"
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/loggers"

	"github.com/gin-gonic/gin"
)


func HandleError(
	c *gin.Context,
	err interface{},
) {

	if appErr, ok := err.(*apperr.Error); ok {
		logErrorMessage(appErr)
		c.JSON(appErr.Status(), appErr.JsonResponse())
	} else {
		unknownAppErr := apperr.New(errors.New("unexpected error"), apperr.UnexpextedError)
		loggers.Errorf("request failed with unknown error: [%+v]", err)
		c.JSON(unknownAppErr.Status(), unknownAppErr.JsonResponse())
	}
}

func logErrorMessage(
	appErr *apperr.Error,
) {
	if appErr.Payload != nil {
		for _, logMessage := range appErr.LogMessages {
			loggers.ErrorWithPayload(
				"failed to make request-id: %s with err: [%v] from ip-address: %v and mac-address: %v;",
				appErr.Payload,
				appErr.RequestID,
				logMessage,
				appErr.IPAddress,
				appErr.MACAddress,
			)
		}
	} else {

		for _, logMessage := range appErr.LogMessages {

			loggers.Errorf(
				"failed to make request: [%s] with err: [%v] from IP [x] [%v] [x] and MAC [x] [%v] [x]",
				appErr.RequestID,
				logMessage,
				appErr.IPAddress,
				appErr.MACAddress,
			)
		}
	}
}