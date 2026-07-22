package server

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ParseServersQuery(c *gin.Context) (GetServersQuery, error) {
	var query GetServersQuery

	// filters
	// if status := c.Query("status"); status != "" {
	// 	parsed, err := strconv.ParseBool(status)
	// 	if err != nil {
	// 		return query, errors.New("Invalid status")
	// 	}
	// 	query.Status = &parsed
	// }

	// if protocol := c.Query("protocol"); protocol != "" {
	// 	query.Protocol = &protocol
	// }

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
		page := int(parsed)
		query.Page = &page
	}
	if pageSize := c.Query("page_size"); pageSize != "" {
		parsed, err := strconv.ParseInt(pageSize, 10, 64)
		if err != nil {
			return query, errors.New("Invalid page size")
		}
		pageSize := int(parsed)
		query.PageSize = &pageSize
	}

	// It's not always about pagination. What if I want to retrieve all for, let's say, export function, right ?
	// So I made it that, if both are nil -> Nothing is parsed in -> No pagination happens
	// If user only passed one param in, I can assume that they want to paginate search result -> Fill the other with default number
	// TODO: What if someone knows this and intentionally spam this request with GET /servers without extra param just to dump everything ? Gotta deal with it
	if query.Page != nil || query.PageSize != nil {
		// Correct values for pages to get. If it's too out of range, fix to the default
		if *query.Page < 1 {
			*query.Page = 1
		}
		if *query.PageSize <= 0 || *query.PageSize > 100 {
			*query.PageSize = 10
		}
	}

	// sorting
	query.SortBy = c.DefaultQuery("sort_by", "id")
	query.Order = c.DefaultQuery("order", "asc")

	return query, nil
}
