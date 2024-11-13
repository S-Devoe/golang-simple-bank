package util

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

const (
	DefaultPage  = 1  // Default page value if not provided
	DefaultLimit = 10 // Default limit value if not provided
)

// PaginationParams holds the pagination parameters like limit, page, and offset.
type PaginationParams struct {
	Limit  int32 // Limit: how many items per page
	Page   int32 // Page: the current page number
	Offset int32 // Offset: the starting index for SQL queries
}

// ParsePaginationQuery is a utility function to parse the pagination query parameters.
// It also applies default values if limit or page is not provided.
func ParsePaginationQuery(ctx *gin.Context) (PaginationParams, error) {
	// Get `limit` and `page` query parameters (defaults if missing)
	limitStr := ctx.DefaultQuery("limit", strconv.Itoa(DefaultLimit)) // Default limit
	pageStr := ctx.DefaultQuery("page", strconv.Itoa(DefaultPage))    // Default page

	// Convert the query parameters to integers
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		return PaginationParams{}, fmt.Errorf("invalid limit: %v", err)
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page <= 0 {
		return PaginationParams{}, fmt.Errorf("invalid page: %v", err)
	}

	// Calculate the offset
	offset := (page - 1) * limit

	// Return pagination parameters
	return PaginationParams{
		Limit:  int32(limit),
		Page:   int32(page),
		Offset: int32(offset),
	}, nil
}
