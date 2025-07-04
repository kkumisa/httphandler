package httphandler

import "net/url"

// Example usage structs demonstrating how to use the package

// UserRequest demonstrates embedding for a user resource
type UserRequest struct {
	IDParam        // Embeds route parameter binding
	Name    string `json:"name"`
	Email   string `json:"email"`
}

// RouteParamName overrides the default to use "user_id" instead of "id"
func (ur *UserRequest) RouteParamName() string {
	return "user_id"
}

// UserListRequest demonstrates query parameter binding
type UserListRequest struct {
	PaginatedList[User]        // Embeds pagination
	SortParams                 // Embeds sorting
	FilterParams               // Embeds filtering
	Status              string `json:"status,omitempty"`
}

// BindQueryParams extends the base implementation with custom parameters
func (ulr *UserListRequest) BindQueryParams(values url.Values) error {
	// First bind pagination params
	if err := ulr.PaginatedList.BindQueryParams(values); err != nil {
		return err
	}

	// Then bind additional custom query params
	ulr.Status = values.Get("status")
	return nil
}

// UserPatchRequest demonstrates PATCH request handling
type UserPatchRequest struct {
	IDParam            // Route parameter binding
	PatchFields        // Field selection for updates
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
}

// User represents a user entity (example)
type User struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
	Age    int    `json:"age"`
}
