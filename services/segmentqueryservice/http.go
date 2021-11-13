package segmentqueryservice

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/usersegment").Subrouter()
	subRouter.HandleFunc("/{segmentUID}", s.list()).Methods("GET")
}

func (s *service) list() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		segmentUID := mux.Vars(r)["segmentUID"]
		users, err := s.GetUsersForSegment(ctx, segmentUID)
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
