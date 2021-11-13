package segmentqueryservice

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func (s *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.HandleFunc("/usersegment/{segmentUID}", s.list()).Methods("GET")
}

func (s *service) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
