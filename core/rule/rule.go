package rule

import (
	"context"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/segment"
)

type UserSegmentRule struct {
	UID         string
	Description string
	UserSegment segment.UserSegmentDefinition
	Action      action.UserActioner
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, name string) (UserSegmentRule, bool, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}
