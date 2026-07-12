package handlers

import (
	"strconv"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"

	"github.com/labstack/echo/v4"
)

func pageableFromQuery(c echo.Context) *dtos.PageableRequest {
	pr := dtos.NewPageableRequest()
	if page, err := strconv.Atoi(c.QueryParam("page")); err == nil && page > 0 {
		pr.Page = page
	}
	if pageSize, err := strconv.Atoi(c.QueryParam("page_size")); err == nil && pageSize > 0 {
		pr.PageSize = pageSize
	}
	dtos.ClampPageSize(pr, constants.MaxPageSize)
	return pr
}
