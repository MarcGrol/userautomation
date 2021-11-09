package segment

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/user"
)

type UserSegment struct {
	UID            string
	Description    string
	UserFilterName string
}

func (us UserSegment) IsApplicableForUser(ctx context.Context, u user.User) (bool, error) {
	filterFunc, found := user.GetUserFilterByName(ctx, us.UserFilterName)
	if !found {
		return false, fmt.Errorf("User filter function with name %s was not found", us.UserFilterName)
	}
	matched, err := filterFunc(ctx, u)
	if err != nil {
		return false, fmt.Errorf("Error filteriing user %s: %s", u.UID, err)
	}
	return matched, nil
}

type UserSegmentManagement interface {
	Put(ctx context.Context, userSegment UserSegment) error
	Get(ctx context.Context, userSegmentUID string) (UserSegment, bool, error)
	List(ctx context.Context) ([]UserSegment, error)
	Remove(ctx context.Context, userUID string) error
}

type Querier interface {
	GetUsersForSegment(ctx context.Context, userSegmentUID string) ([]user.User, error)
}
