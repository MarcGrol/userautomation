package userruletriggerservice

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/segmentrule").Subrouter()
	subRouter.HandleFunc("", m.post()).Methods("POST")
}

func (m *service) post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
