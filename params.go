package httphandler

import (
	"net/url"
	"strings"
)

// IDParam provides route parameter binding for ID-based resources
// This Object satisfies RoutParamBinder interface and can be embedded into object that needs
// route parameter binding.
type IDParam struct {
	ID string `json:"id"`
}

// RouteParamName returns the route parameter name for the ID.
// Override this in embedded structs if route param name should be different than 'id'.
func (p *IDParam) RouteParamName() string {
	return "id"
}

func (p *IDParam) BindRouteParam(value string) {
	p.ID = value
}

// PatchFields handles field selection for PATCH operations.
// This Object satisfies PatchFieldExtractor interface and can be embedded into object that needs
// extraction of the patch field.
type PatchFields struct {
	fieldsToUpdate []string
}

func (pf *PatchFields) ExtractPatchFields(values url.Values) error {
	fieldsParam := values.Get("fields")
	if fieldsParam == "" {
		return NewBadRequestError("'fields' query parameter is required for PATCH requests")
	}

	fields := strings.Split(fieldsParam, ",")

	// Clean and validate fields
	var cleanFields []string
	for _, field := range fields {
		if trimmed := strings.TrimSpace(field); trimmed != "" {
			cleanFields = append(cleanFields, trimmed)
		}
	}

	if len(cleanFields) == 0 {
		return NewBadRequestError("at least one field must be specified for update")
	}

	pf.fieldsToUpdate = cleanFields
	return nil
}
