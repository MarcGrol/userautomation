package actionmanager

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/action").Subrouter()
	subRouter.HandleFunc("", m.listActionSpecs()).Methods("GET")

}

func (m *service) listActionSpecs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
