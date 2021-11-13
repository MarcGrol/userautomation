package actionmanager

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/action").Subrouter()
	subRouter.HandleFunc("", s.listActionSpecs()).Methods("GET")

}

func (s *service) listActionSpecs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		actions, err := s.ListActionSpecs(ctx)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(actions)
		w.WriteHeader(http.StatusOK)
	}
}
