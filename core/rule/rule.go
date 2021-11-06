package rule

import (
	"context"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/segment"
)

type TriggerKind int

const (
	TriggerOnDemand TriggerKind = 1 << iota
	TriggerCron
	TriggerUserChange
)

type UserSegmentRule struct {
	UID             string
	Description     string
	UserSegment     segment.UserSegment
	Action          action.UserActioner
	TriggerKindMask TriggerKind
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, ruleUID string) (UserSegmentRule, bool, error)
	Delete(ctx context.Context, ruleUID string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}

type SegmentRuleExecutionService interface {
	Execute(ctx context.Context, ruleUID string) error
}
