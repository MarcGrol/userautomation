package util

import (
	"context"
)

type PreProvisioner interface {
	Preprov(ctx context.Context) error
}

type NoPreProvNeeded struct{}

func (_ NoPreProvNeeded) Preprov(ctx context.Context) error {
	return nil
}
