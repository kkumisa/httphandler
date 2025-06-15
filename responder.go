package httphandler

import (
	"encoding/json"
	"errors"
	"net/http"
)

// HTTPResponder handles writing HTTP responses.
type HTTPResponder struct{}

func NewHTTPResponder() *HTTPResponder {
	return &HTTPResponder{}
}

// RespondWithError writes an error response to the HTTP response writer.
func (r *HTTPResponder) RespondWithError(w http.ResponseWriter, err error) {
	var httpErr HTTPError
	if !errors.As(err, &httpErr) {
		httpErr = NewInternalServerError("Internal server error")
	}

	response := ErrorResponse{Error: httpErr.Message()}

	jsonData, encodeErr := json.Marshal(response)
	if encodeErr != nil {
		//TODO see what to do with this
		r.respondWithPlainTextError(w, httpErr)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpErr.StatusCode())
	w.Write(jsonData)
}

// RespondWithSuccess writes a successful response to the HTTP response writer.
func (r *HTTPResponder) RespondWithSuccess(w http.ResponseWriter, data any) {
	if data == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		r.RespondWithError(w, NewInternalServerError("failed to encode response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

// respondWithPlainTextError writes error as plain text when JSON encoding fails.
func (r *HTTPResponder) respondWithPlainTextError(w http.ResponseWriter, err HTTPError) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(err.StatusCode())
	w.Write([]byte(err.Message()))
}
