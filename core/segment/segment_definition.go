package segment

import (
	"context"

	"github.com/MarcGrol/userautomation/core/user"
)

type UserSegmentDefinition struct {
	UID                 string
	Description         string
	IsApplicableForUser user.UserFilterFunc // Could use a WHERE clause alternatively
}

type UserSegmentService interface {
	Put(ctx context.Context, userSegment UserSegmentDefinition) error
	Get(ctx context.Context, userSegmentUID string) (UserSegmentDefinition, bool, error)
	List(ctx context.Context) ([]UserSegmentDefinition, error)
	Remove(ctx context.Context, userUID string) error
	GetUsersForSegment(ctx context.Context, userSegmentUID string) ([]user.User, error)
}
