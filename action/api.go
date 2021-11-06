package action

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/users"
)

type UserChangeType int

const (
	UserCreated UserChangeType = iota
	UserModified
	UserRemoved
)

type UserActioner interface {
	Perform(ctx context.Context, action UserAction) error
}

type UserAction struct {
	RuleName       string
	UserChangeType UserChangeType
	OldState       *users.User
	NewState       *users.User
}

func (a UserAction) String() string {
	return fmt.Sprintf("UserActioner for rule '%s' triggered action om User '%s' - status: %+v\n",
		a.RuleName, getUserUid(a.OldState, a.NewState), a.UserChangeType)
}

func getUserUid(oldState *users.User, newState *users.User) string {
	if oldState != nil {
		return oldState.UID
	}
	return newState.UID
}
