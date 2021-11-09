package rule

import (
	"context"

	"github.com/MarcGrol/userautomation/core/segment"
)

type TriggerAllowed int

const (
	TriggerOnDemand TriggerAllowed = 1 << iota
	TriggerCron
	TriggerUserChange
)

type RuleSpec struct {
	UID             string
	Description     string
	SegmentSpec     segment.SegmentSpec
	//Task          usertask.UserTaskExecutor
	ActionName string
	AllowedTriggers TriggerAllowed
}

type RuleService interface {
	Put(ctx context.Context, segmentRule RuleSpec) error
	Get(ctx context.Context, ruleUID string) (RuleSpec, bool, error)
	Remove(ctx context.Context, ruleUID string) error
	List(ctx context.Context) ([]RuleSpec, error)
}

//go:generate mockgen -source=rule.go -destination=rule_execution_mock.go -package=rule TriggerRuleExecution
type TriggerRuleExecution interface {
	Trigger(ctx context.Context, ruleUID string) error
}
