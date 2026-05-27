package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseServersQuery(c *gin.Context) (GetServersQuery, error) {
	var query GetServersQuery

	// filters
	if status := c.Query("status"); status != "" {
		parsed, err := strconv.ParseBool(status)
		if err != nil {
			return query, errors.New("Invalid status")
		}
		query.Status = &parsed
	}

	if protocol := c.Query("protocol"); protocol != "" {
		query.Protocol = &protocol
	}

	if name := c.Query("name"); name != "" {
		query.Name = &name
	}

	// Pagination
	// Validate input
	if page := c.Query("page"); page != "" {
		parsed, err := strconv.ParseInt(page, 10, 64)
		if err != nil {
			return query, errors.New("Invalid page")
		}
		query.Page = int(parsed)
	}
	if pageSize := c.Query("page"); pageSize != "" {
		parsed, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			return query, errors.New("Invalid page size")
		}
		query.PageSize = int(parsed)
	}

	// Correct values for pages to get. If it's too out of range, fix to the default
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize <= 0 || query.PageSize > 100 {
		query.PageSize = 10
	}

	// sorting
	query.SortBy = c.DefaultQuery("sort_by", "id")
	query.Order = c.DefaultQuery("order", "asc")

	return query, nil
}
