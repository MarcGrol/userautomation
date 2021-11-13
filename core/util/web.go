package util

import (
	"context"

	"github.com/gorilla/mux"
)

type WebExposer interface {
	RegisterEndpoints(ctx context.Context, router *mux.Router)
}

type NoWebNeeded struct{}

func (_ NoWebNeeded) RegisterEndpoints(ctx context.Context, router *mux.Router) {}
