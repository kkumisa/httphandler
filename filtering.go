package httphandler

import (
	"fmt"
	"net/url"
	"strings"
)

type FilterOperator string

const (
	FilterEq       FilterOperator = "eq"       // Equal
	FilterNe       FilterOperator = "ne"       // Not equal
	FilterGt       FilterOperator = "gt"       // Greater than
	FilterGte      FilterOperator = "gte"      // Greater than or equal
	FilterLt       FilterOperator = "lt"       // Less than
	FilterLte      FilterOperator = "lte"      // Less than or equal
	FilterIn       FilterOperator = "in"       // In list
	FilterNotIn    FilterOperator = "not_in"   // Not in list
	FilterContains FilterOperator = "contains" // Contains substring
	FilterPrefix   FilterOperator = "prefix"   // Starts with
	FilterSuffix   FilterOperator = "suffix"   // Ends with
	FilterIsNull   FilterOperator = "is_null"  // Is null
	FilterNotNull  FilterOperator = "not_null" // Is not null
)

// FilterCondition represents a single filter condition
type FilterCondition struct {
	Field    string         `json:"field"`
	Operator FilterOperator `json:"operator"`
	Value    string         `json:"value"`
	Values   []string       `json:"values,omitempty"` // For 'in' and 'not_in' operators
}

type FilterParams struct {
	Conditions []FilterCondition `json:"-"`
}

// BindFilters implements FilterBinder interface
func (fp *FilterParams) BindFilters(values url.Values) error {
	fp.Conditions = []FilterCondition{}

	for key, paramValues := range values {
		// Skip non-filter parameters
		if !strings.HasPrefix(key, "filter.") {
			continue
		}

		// Extract field and operator from key
		// Format: filter.field_name[operator] or filter.field_name (defaults to 'eq')
		filterKey := key[7:] // Remove "filter." prefix

		field, operator, err := fp.parseFilterKey(filterKey)
		if err != nil {
			return NewBadRequestError(fmt.Sprintf("invalid filter key '%s': %v", key, err))
		}

		if len(paramValues) == 0 {
			continue
		}

		condition := FilterCondition{
			Field:    field,
			Operator: operator,
		}

		// Handle operators that don't need values
		if operator == FilterIsNull || operator == FilterNotNull {
			fp.Conditions = append(fp.Conditions, condition)
			continue
		}

		// Handle multi-value operators
		if operator == FilterIn || operator == FilterNotIn {
			// Split comma-separated values or use multiple query parameters
			var allValues []string
			for _, paramValue := range paramValues {
				values := strings.Split(paramValue, ",")
				for _, v := range values {
					if trimmed := strings.TrimSpace(v); trimmed != "" {
						allValues = append(allValues, trimmed)
					}
				}
			}
			if len(allValues) == 0 {
				return NewBadRequestError(fmt.Sprintf("filter '%s' requires at least one value", key))
			}
			condition.Values = allValues
		} else {
			// Single value operators
			condition.Value = strings.TrimSpace(paramValues[0])
			if condition.Value == "" {
				return NewBadRequestError(fmt.Sprintf("filter '%s' requires a value", key))
			}
		}

		fp.Conditions = append(fp.Conditions, condition)
	}

	return nil
}

func (fp *FilterParams) parseFilterKey(filterKey string) (field string, operator FilterOperator, err error) {
	// Default operator
	operator = FilterEq

	// Check for operator in brackets: field_name[operator]
	if strings.Contains(filterKey, "[") && strings.HasSuffix(filterKey, "]") {
		parts := strings.SplitN(filterKey, "[", 2)
		if len(parts) != 2 {
			return "", "", fmt.Errorf("invalid filter format")
		}

		field = parts[0]
		operatorStr := strings.TrimSuffix(parts[1], "]")

		operator = FilterOperator(operatorStr)
		if !fp.isValidOperator(operator) {
			return "", "", fmt.Errorf("unsupported operator '%s'", operatorStr)
		}
	} else {
		field = filterKey
	}

	if field == "" {
		return "", "", fmt.Errorf("field name cannot be empty")
	}

	return field, operator, nil
}

func (fp *FilterParams) isValidOperator(op FilterOperator) bool {
	switch op {
	case FilterEq, FilterNe, FilterGt, FilterGte, FilterLt, FilterLte,
		FilterIn, FilterNotIn, FilterContains, FilterPrefix, FilterSuffix,
		FilterIsNull, FilterNotNull:
		return true
	default:
		return false
	}
}
