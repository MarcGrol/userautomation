package userruletriggerservice

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/MarcGrol/userautomation/core/userrule"

	"github.com/gorilla/mux"
)

func (s *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {
	subRouter := router.PathPrefix("/api/userrule").Subrouter()
	subRouter.HandleFunc("", s.post()).Methods("POST")
}

func (s *service) post() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		rule := userrule.Spec{}
		err := json.NewDecoder(r.Body).Decode(&rule)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		err = s.Trigger(ctx, rule)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(err.Error())
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
