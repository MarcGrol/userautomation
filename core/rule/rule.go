package rule

import (
	"context"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/segment"
)

type TriggerAllowed int

const (
	TriggerOnDemand TriggerAllowed = 1 << iota
	TriggerCron
	TriggerUserChange
)

type UserSegmentRule struct {
	UID             string
	Description     string
	UserSegment     segment.UserSegment
	Action          action.UserActioner
	AllowedTriggers TriggerAllowed
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, ruleUID string) (UserSegmentRule, bool, error)
	Remove(ctx context.Context, ruleUID string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}

//go:generate mockgen -source=rule.go -destination=rule_execution_mock.go -package=rule SegmentRuleExecutionService
type SegmentRuleExecutionService interface {
	Trigger(ctx context.Context, ruleUID string) error
}
