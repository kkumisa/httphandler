package httphandler

import "net/url"

// Request behavior interfaces define how different types of requests should be processed.

// RouteParamBinder extracts route parameters and binds them to the request.
type RouteParamBinder interface {
	RouteParamName() string
	BindRouteParam(value string)
}

// QueryParamBinder extracts query parameters and binds them to the request.
type QueryParamBinder interface {
	BindQueryParams(values url.Values) error
}

// FilterBinder extracts and validates filter parameters from query string.
type FilterBinder interface {
	BindFilters(values url.Values) error
}

// SortBinder extracts and validates sort parameters from query string.
type SortBinder interface {
	BindSort(values url.Values) error
}

// PatchFieldExtractor extracts and validates fields to be updated in PATCH requests.
type PatchFieldExtractor interface {
	ExtractPatchFields(values url.Values) error
}
