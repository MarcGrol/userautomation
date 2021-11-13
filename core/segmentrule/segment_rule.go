package segmentrule

import (
	"context"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/util"
)

type Spec struct {
	UID         string
	Description string
	SegmentSpec segment.Spec
	ActionSpec  action.Spec
}

type Management interface {
	Put(ctx context.Context, rule Spec) error
	Get(ctx context.Context, ruleUID string) (Spec, bool, error)
	Remove(ctx context.Context, ruleUID string) error
	List(ctx context.Context) ([]Spec, error)
	util.PreProvisioner
	util.WebExposer
}

type TriggerRuleExecution interface {
	Trigger(ctx context.Context, rule Spec) error
	util.WebExposer
}
