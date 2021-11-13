package actionmanager

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.HandleFunc("/action", m.listActionSpecs()).Methods("GET")

}

func (m *service) listActionSpecs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
