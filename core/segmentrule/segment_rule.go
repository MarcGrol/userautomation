package segmentrule

import (
	"context"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/segment"
)

type Spec struct {
	UID         string
	Description string
	SegmentSpec segment.Spec
	ActionSpec  action.Spec
}

type Service interface {
	Put(ctx context.Context, rule Spec) error
	Get(ctx context.Context, ruleUID string) (Spec, bool, error)
	Remove(ctx context.Context, ruleUID string) error
	List(ctx context.Context) ([]Spec, error)
}

//go:generate mockgen -source=segment_rule.go -destination=rule_execution_mock.go -package=segmentrule TriggerRuleExecution
type TriggerRuleExecution interface {
	Trigger(ctx context.Context, rule Spec) error
}
