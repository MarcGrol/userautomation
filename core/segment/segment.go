package segment

import (
	"context"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/core/util"
)

type Spec struct {
	UID            string
	Description    string
	UserFilterName string
}

type Management interface {
	Put(ctx context.Context, segment Spec) error
	Get(ctx context.Context, segmentUID string) (Spec, bool, error)
	List(ctx context.Context) ([]Spec, error)
	Remove(ctx context.Context, segmentUID string) error
	util.PreProvisioner
	util.WebExposer
}

type Querier interface {
	GetUsersForSegment(ctx context.Context, segmentUID string) ([]user.User, error)
	util.WebExposer
}
