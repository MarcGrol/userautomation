package action

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/user"
)

type UserActioner interface {
	Perform(ctx context.Context, action UserAction) error
}

type UserChangeType int

const (
	UserCreated UserChangeType = iota
	UserModified
	UserRemoved
)

type UserAction struct {
	RuleName       string
	UserChangeType UserChangeType
	OldState       *user.User
	NewState       *user.User
}

func (a UserAction) String() string {
	return fmt.Sprintf("UserActioner for rule '%s' triggered action om User '%s' - status: %+v\n",
		a.RuleName, getUserUID(a.OldState, a.NewState), a.UserChangeType)
}

func getUserUID(oldState *user.User, newState *user.User) string {
	if oldState != nil {
		return oldState.UID
	}
	return newState.UID
}
