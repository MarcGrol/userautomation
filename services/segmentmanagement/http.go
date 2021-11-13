package segmentmanagement

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
)

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/segment").Subrouter()
	subRouter.HandleFunc("", m.list()).Methods("GET")
	subRouter.HandleFunc("/{ruleUID}", m.get()).Methods("GET")
	subRouter.HandleFunc("/{ruleUID}", m.put()).Methods("PUT")
	subRouter.HandleFunc("/{ruleUID}", m.remove()).Methods("DELETE")
}

func (m *service) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (m *service) put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (m *service) remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (m *service) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}
