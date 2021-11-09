package action

import (
	"context"
)

type ActionSpec struct {
	Name                    string
	Description             string
	MandatoryUserAttributes []string
	ProvidedAttributes      map[string]string
}

type ActionManager interface {
	GetActionSpecOnName(ctx context.Context, name string) (ActionSpec, bool, error)
	ListActionSpecs(ctx context.Context) ([]ActionSpec, error)
}
