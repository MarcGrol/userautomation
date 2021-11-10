package segment

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/user"
)

type Spec struct {
	UID            string
	Description    string
	UserFilterName string
}

func (s Spec) IsApplicableForUser(ctx context.Context, u user.User) (bool, error) {
	filterFunc, found := user.GetUserFilterByName(ctx, s.UserFilterName)
	if !found {
		return false, fmt.Errorf("User filter function with name %s was not found", s.UserFilterName)
	}
	matched, err := filterFunc(ctx, u)
	if err != nil {
		return false, fmt.Errorf("Error filteriing user %s: %s", u.UID, err)
	}
	return matched, nil
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
