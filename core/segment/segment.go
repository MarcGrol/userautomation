package segment

import (
	"context"

	"github.com/MarcGrol/userautomation/core/user"
)

type UserSegment struct {
	UID         string
	Description string
	//IsApplicableForUser user.FilterFunc // Could use a WHERE clause alternatively
	UserFilterName string
	Users          map[string]user.User
}

type UserSegmentService interface {
	Put(ctx context.Context, userSegment UserSegment) error
	Get(ctx context.Context, userSegmentUID string) (UserSegment, bool, error)
	List(ctx context.Context) ([]UserSegment, error)
	Remove(ctx context.Context, userUID string) error
}

type UserSegmentQueryService interface {
	GetUsersForSegment(ctx context.Context, userSegmentUID string) ([]user.User, error)
}
