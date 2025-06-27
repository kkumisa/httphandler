package httphandler

import (
	"context"
	"net/http"
)

// Handler represents a generic HTTP handler function signature
type Handler[TReq any, TRes any] func(ctx context.Context, req *TReq) (*TRes, error)

// NewGenericHandler creates HTTP middleware for generic request handling
func NewGenericHandler[TReq any, TRes any](handler Handler[TReq, TRes]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TReq

		if err := BindRequest(r, &req); err != nil {
			RespondWithError(w, err)
			return
		}

		ctx := r.Context()
		response, err := handler(ctx, &req)
		if err != nil {
			RespondWithError(w, err)
			return
		}

		// TODO Add support for other 200 family statuses
		RespondWithSuccess(w, response)
	}
}
