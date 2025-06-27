package httphandler

import (
	"net/url"
	"strings"
)

type SortDirection string

const (
	SortAsc  SortDirection = "asc"
	SortDesc SortDirection = "desc"
)

// SortField represents a field to sort by with its direction
type SortField struct {
	Field     string        `json:"field"`
	Direction SortDirection `json:"direction"`
}

type SortParams struct {
	Fields []SortField `json:"-"`
}

// BindSort implements SortBinder interface
func (sp *SortParams) BindSort(values url.Values) error {
	sortParam := values.Get("sort")
	if sortParam == "" {
		// Sorting is optional
		return nil
	}

	fields := strings.Split(sortParam, ",")
	sp.Fields = make([]SortField, 0, len(fields))

	for _, field := range fields {
		field = strings.TrimSpace(field)
		if field == "" {
			continue
		}

		sortField := SortField{
			// Default direction
			Direction: SortAsc,
		}

		switch {
		case strings.HasPrefix(field, "-"):
			sortField.Direction = SortDesc
			sortField.Field = field[1:]
		case strings.HasPrefix(field, "+"):
			sortField.Direction = SortAsc
			sortField.Field = field[1:]
		default:
			sortField.Field = field
		}

		sp.Fields = append(sp.Fields, sortField)
	}

	return nil
}
