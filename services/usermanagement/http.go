package usermanagement

import (
	"context"
	"github.com/gorilla/mux"
	"net/http"
)

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api").Subrouter()
	subRouter.HandleFunc("/user", m.list()).Methods("GET")
	subRouter.HandleFunc("/user/{ruleUID}", m.get()).Methods("GET")
	subRouter.HandleFunc("/user/{ruleUID}", m.put()).Methods("PUT")
	subRouter.HandleFunc("/user/{ruleUID}", m.remove()).Methods("DELETE")
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
