package segment

import (
	"context"
	"github.com/MarcGrol/userautomation/core/user"
)

type Spec struct {
	UID            string
	Description    string
	UserFilterName string
}

type SegmentWithUsers struct {
	SegmentSpec Spec
	Users       map[string]user.User
}

type Management interface {
	Put(ctx context.Context, segment Spec) error
	Get(ctx context.Context, segmentUID string) (Spec, bool, error)
	List(ctx context.Context) ([]Spec, error)
	Remove(ctx context.Context, segmentUID string) error
}

type Querier interface {
	GetUsersForSegment(ctx context.Context, segmentUID string) ([]user.User, error)
}
