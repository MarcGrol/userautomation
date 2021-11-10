package action

import (
	"context"
)

type Spec struct {
	Name                    string
	Description             string
	MandatoryUserAttributes []string
	ProvidedInformation     map[string]string
}

type ActionManager interface {
	GetActionSpecOnName(ctx context.Context, name string) (Spec, bool, error)
	ListActionSpecs(ctx context.Context) ([]Spec, error)
}
