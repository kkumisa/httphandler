package httphandler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// BindRequest processes an HTTP request and binds data to the request struct.
func BindRequest(r *http.Request, req any) error {
	// Bind route parameters (e.g., /users/{user_id})
	if err := bindRouteParams(r, req); err != nil {
		return fmt.Errorf("couldn't bind route param: %w", err)
	}

	// Method-specific binding
	switch r.Method {
	case http.MethodGet:
		if err := bindQueryParams(r, req); err != nil {
			return fmt.Errorf("couldn't bind query param: %w", err)
		}
	case http.MethodPatch:
		if err := bindPatchRequest(r, req); err != nil {
			return fmt.Errorf("couldn't bind patch request: %w ", err)
		}

	default:
		if err := bindJSONBody(r, req); err != nil {
			return fmt.Errorf("couldn't bind json body: %w ", err)
		}
	}

	return nil
}

// bindRouteParams extracts route parameters and binds them to the request.
func bindRouteParams(r *http.Request, req any) error {
	binder, ok := req.(RouteParamBinder)
	if !ok {
		return nil
	}

	paramName := binder.RouteParamName()
	value := r.PathValue(paramName)

	if value == "" {
		return NewBadRequestError(fmt.Sprintf("missing required route parameter: %s", paramName))
	}

	binder.BindRouteParam(value)
	return nil
}

// bindQueryParams processes query parameters for GET requests.
func bindQueryParams(r *http.Request, req any) error {
	queryParams := r.URL.Query()

	if binder, ok := req.(QueryParamBinder); ok {
		if err := binder.BindQueryParams(queryParams); err != nil {
			return err
		}
	}

	if filterBinder, ok := req.(FilterBinder); ok {
		if err := filterBinder.BindFilters(queryParams); err != nil {
			return err
		}
	}

	if sortBinder, ok := req.(SortBinder); ok {
		if err := sortBinder.BindSort(queryParams); err != nil {
			return err
		}
	}

	return nil
}

// bindPatchRequest handles PATCH request binding (both query params and JSON body).
func bindPatchRequest(r *http.Request, req any) error {
	extractor, ok := req.(PatchFieldExtractor)
	if !ok {
		return bindJSONBody(r, req)
	}

	if err := extractor.ExtractPatchFields(r.URL.Query()); err != nil {
		return fmt.Errorf("failed to extract patch fields: %w", err)
	}

	return bindJSONBody(r, req)
}

// bindJSONBody decodes JSON request body into the request struct.
func bindJSONBody(r *http.Request, req any) error {
	if r.Body == nil {
		return NewBadRequestError("request body is required")
	}

	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		return NewBadRequestError("invalid JSON in request body")
	}

	return nil
}
