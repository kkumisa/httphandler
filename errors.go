package httphandler

import "net/http"

type HTTPError struct {
	code    int
	message string
}

func NewBadRequestError(message string) HTTPError {
	return HTTPError{code: http.StatusBadRequest, message: message}
}

func NewInternalServerError(message string) HTTPError {
	return HTTPError{code: http.StatusInternalServerError, message: message}
}

func NewNotFoundError(message string) HTTPError {
	return HTTPError{code: http.StatusNotFound, message: message}
}

func NewUnauthorizedError(message string) HTTPError {
	return HTTPError{code: http.StatusUnauthorized, message: message}
}

func NewForbiddenError(message string) HTTPError {
	return HTTPError{code: http.StatusForbidden, message: message}
}

func NewConflictError(message string) HTTPError {
	return HTTPError{code: http.StatusConflict, message: message}
}

func NewGenericError(message string, code int) HTTPError {
	return HTTPError{code: code, message: message}
}

func (e HTTPError) Error() string {
	return e.message
}

func (e HTTPError) StatusCode() int {
	return e.code
}

func (e HTTPError) Message() string {
	return e.message
}

type ErrorResponse struct {
	Error string `json:"error"`
}
