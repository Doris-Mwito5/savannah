package ctxfilter

import (
	"github/Doris-Mwito5/savannah-pos/internal/apperr"
	"github/Doris-Mwito5/savannah-pos/internal/models"
	"github/Doris-Mwito5/savannah-pos/internal/null"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func FilterFromContext(
	c *gin.Context,
) (*models.Filter, error) {

	filter := &models.Filter{}

	page, per, err := paginationFromContext(c)
	if err != nil {
		return filter, err
	}

	filter.Page = page
	filter.Per = per
	filter.From = strings.TrimSpace(c.Query("from"))
	filter.To = strings.TrimSpace(c.Query("to"))
	filter.Status = null.NullValue(strings.TrimSpace(c.Query("status")))
	filter.Type = strings.TrimSpace(c.Query("type"))
	filter.Token = strings.TrimSpace(c.Query("token"))
	filter.Term = strings.TrimSpace(c.Query("term"))
	filter.UUID = strings.TrimSpace(c.Query("uuid"))
	filter.Year = strings.TrimSpace(c.Query("year"))
	filter.Reference = strings.TrimSpace(c.Query("reference"))

	return filter, nil
}

func paginationFromContext(
	c *gin.Context,
) (int, int, error) {

	page := 1
	per := 20

	var err error

	pageQuery := strings.TrimSpace(c.Query("page"))
	if pageQuery != "" {
		page, err = strconv.Atoi(pageQuery)
		if err != nil {
			log.Printf("Failed to parse page query param [%v]", pageQuery)
			return page, per, apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
		}
	}

	perQuery := strings.TrimSpace(c.Query("per"))
	if perQuery != "" {
		per, err = strconv.Atoi(perQuery)
		if err != nil {
			log.Printf("Failed to parse per query param [%v]", perQuery)
			return page, per, apperr.NewErrorWithType(
				err,
				apperr.BadRequest,
			)
		}
	}

	return page, per, nil
}
