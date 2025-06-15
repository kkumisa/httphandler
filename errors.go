package httphandler

import "net/http"

// HTTPError represents an HTTP error with status code and message
type HTTPError struct {
	code    int
	message string
}

//TODO Add more of these errors

// NewBadRequestError creates a 400 Bad Request error
func NewBadRequestError(message string) HTTPError {
	return HTTPError{code: http.StatusBadRequest, message: message}
}

// NewInternalServerError creates a 500 Internal Server Error
func NewInternalServerError(message string) HTTPError {
	return HTTPError{code: http.StatusInternalServerError, message: message}
}

// NewNotFoundError creates a 404 Not Found error
func NewNotFoundError(message string) HTTPError {
	return HTTPError{code: http.StatusNotFound, message: message}
}

// NewUnauthorizedError creates a 401 Unauthorized error
func NewUnauthorizedError(message string) HTTPError {
	return HTTPError{code: http.StatusUnauthorized, message: message}
}

// NewForbiddenError creates a 403 Forbidden error
func NewForbiddenError(message string) HTTPError {
	return HTTPError{code: http.StatusForbidden, message: message}
}

// NewConflictError creates a 409 Conflict error
func NewConflictError(message string) HTTPError {
	return HTTPError{code: http.StatusConflict, message: message}
}

// NewGenericError creates an error with a given code
func NewGenericError(message string, code int) HTTPError {
	return HTTPError{code: code, message: message}
}

// Error implements the error interface
func (e HTTPError) Error() string {
	return e.message
}

// StatusCode returns the HTTP status code
func (e HTTPError) StatusCode() int {
	return e.code
}

// Message returns the error message
func (e HTTPError) Message() string {
	return e.message
}

// ErrorResponse represents the JSON structure for error responses
type ErrorResponse struct {
	Error string `json:"error"`
}
