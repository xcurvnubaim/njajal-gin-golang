package query

import (
	"strconv"

	"github.com/gin-gonic/gin"
	CustomValidator "github.com/xcurvnubaim/njajal-gin-golang/internal/pkg/validator"
)

type QueryParams struct {
	Search        *string
	SearchColumns []string
	Filters       *map[string]string
	OrderBy       string
	OrderDir      string
	Page          int
	PageSize      int
}

func NewQueryParams(searchColumns []string) *QueryParams {
	return &QueryParams{SearchColumns: searchColumns}
}

func (qp *QueryParams) Parse(c *gin.Context, defaultPageSize string) {
	search := c.Query("search")
	if search != "" {
		qp.Search = &search
	}

	filters := make(map[string]string)
	for key, value := range c.Request.URL.Query() {
		if key != "search" && key != "order_by" && key != "order_dir" && key != "page" && key != "page_size" {
			filters[key] = value[0]
		}
	}

	if len(filters) > 0 {
		qp.Filters = &filters
	}

	qp.OrderBy = c.DefaultQuery("order_by", "created_at")
	qp.OrderDir = c.DefaultQuery("order_dir", "asc")

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		page = 1
	}
	qp.Page = page

	pageSize, err := strconv.Atoi(c.DefaultQuery("page_size", defaultPageSize))
	if err != nil {
		pageSize = 10
	}
	qp.PageSize = pageSize
}

func (qp *QueryParams) Validate(validator CustomValidator.ParamValidator) error {
	if err := CustomValidator.ValidateSearch(qp.Search, validator.MaxSearchLength); err != nil {
		return err
	}

	if err := CustomValidator.ValidateFilters(qp.Filters, validator.AllowedFilterKeys, validator.MaxFilterValueLength); err != nil {
		return err
	}

	if err := CustomValidator.ValidateOrderBy(qp.OrderBy, validator.AllowedOrderByColumns); err != nil {
		return err
	}

	if err := CustomValidator.ValidateOrderDir(qp.OrderDir); err != nil {
		return err
	}

	if err := CustomValidator.ValidatePage(qp.Page, qp.PageSize, validator.MaxPageSize); err != nil {
		return err
	}

	return nil
}
