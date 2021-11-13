package util

import (
	"context"
	"github.com/gorilla/mux"
)

type PreProvisioner interface {
	Preprov(ctx context.Context) error
}

type WebExposer interface {
	RegisterEndpoints(ctx context.Context, router *mux.Router)
}

type NoPreProvNeeded struct{}

func (_ NoPreProvNeeded) Preprov(ctx context.Context) error {
	return nil
}

type NoWebNeeded struct{}

func (_ NoWebNeeded) RegisterEndpoints(ctx context.Context, router *mux.Router) {}
