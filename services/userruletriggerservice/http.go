package userruletriggerservice

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.HandleFunc("/segmentrule", m.post()).Methods("POST")
}

func (m *service) post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
