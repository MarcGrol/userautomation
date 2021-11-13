package usermanagement

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/user").Subrouter()
	subRouter.HandleFunc("", s.list()).Methods("GET")
	subRouter.HandleFunc("/{ruleUID}", s.get()).Methods("GET")
	subRouter.HandleFunc("/{ruleUID}", s.put()).Methods("PUT")
	subRouter.HandleFunc("/{ruleUID}", s.remove()).Methods("DELETE")
}

func (s *service) get() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (s *service) put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (s *service) remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotImplemented)
	}
}

func (m *service) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		users, err := m.List(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}
