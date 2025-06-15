package httphandler

import (
	"net/url"
	"strconv"
)

// PaginatedList represents a paginated collection of items
type PaginatedList[T any] struct {
	Limit      int    `json:"limit"`
	NextCursor string `json:"next_cursor,omitempty"`
	Items      []T    `json:"items"`
}

// BindQueryParams implements QueryParamBinder interface
func (pl *PaginatedList[T]) BindQueryParams(values url.Values) error {
	pl.setDefaults()

	if limitStr := values.Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil || limit <= 0 {
			return NewBadRequestError("limit must be a positive integer")
		}

		// Prevent abuse with reasonable upper limit
		if limit > 1000 {
			limit = 1000
		}

		pl.Limit = limit
	}

	pl.NextCursor = values.Get("next_cursor")
	return nil
}

// setDefaults initializes default values
func (pl *PaginatedList[T]) setDefaults() {
	pl.Limit = 10
	pl.NextCursor = ""
	if pl.Items == nil {
		pl.Items = make([]T, 0)
	}
}
