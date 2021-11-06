package rules

import (
	"context"
	"github.com/MarcGrol/userautomation/segments"
)

type UserSegmentRule struct {
	Name        string
	UserSegment segments.UserSegment
	//IsApplicableForUser  users.UserFilterFunc // Could use a WHERE clause alternatively
	PerformActionForUser UserActionFunc
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, name string) (UserSegmentRule, bool, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}
