package rules

import (
	"context"
	"github.com/MarcGrol/userautomation/action"
	"github.com/MarcGrol/userautomation/segments"
)

type UserSegmentRule struct {
	Name        string
	UserSegment segments.UserSegment
	Action      action.UserActioner
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, name string) (UserSegmentRule, bool, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}
