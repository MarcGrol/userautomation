package usermanagement

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MarcGrol/userautomation/core/user"
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
		ctx := context.Background()

		ruleUID := mux.Vars(r)["ruleUID"]
		user, exists, err := s.Get(ctx, ruleUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
		w.WriteHeader(http.StatusOK)
	}
}

func (s *service) put() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		u := user.User{}
		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		err = s.Put(ctx, u)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *service) remove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		ruleUID := mux.Vars(r)["ruleUID"]
		err := s.Remove(ctx, ruleUID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func (s *service) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		users, err := s.List(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
		w.WriteHeader(http.StatusOK)
	}
}
