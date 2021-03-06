package action

import (
	"context"

	"github.com/MarcGrol/userautomation/core/util"
)

type Spec struct {
	Name                    string
	Description             string
	MandatoryUserAttributes []string
	ProvidedInformation     map[string]string
}

type Management interface {
	GetActionSpecOnName(ctx context.Context, name string) (Spec, bool, error)
	ListActionSpecs(ctx context.Context) ([]Spec, error)
	util.PreProvisioner
	util.WebExposer
}
