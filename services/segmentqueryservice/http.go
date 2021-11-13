package segmentqueryservice

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/usersegment").Subrouter()
	subRouter.HandleFunc("/{segmentUID}", s.list()).Methods("GET")
}

func (s *service) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
