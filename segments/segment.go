package segments

import (
	"context"
	"github.com/MarcGrol/userautomation/users"
)

type UserSegment struct {
	Name                string
	IsApplicableForUser users.UserFilterFunc
	Users               []users.User
}

type UserSegmentService interface {
	Put(ctx context.Context, userSegment UserSegment) error
	Get(ctx context.Context, userSegmentUID string) (UserSegment, bool, error)
	List(ctx context.Context) ([]UserSegment, error)
	Remove(ctx context.Context, userUID string) error
}
