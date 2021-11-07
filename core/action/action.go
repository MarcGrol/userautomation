package action

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/user"
)

//go:generate mockgen -source=action.go -destination=actioner_mock.go -package=action UserActioner
type UserActioner interface {
	Perform(ctx context.Context, action UserAction) error
}

type TriggerType int

const (
	UserCreated TriggerType = iota
	UserModified
	UserRemoved
	OnDemand
	Cron
)

type UserAction struct {
	RuleName    string
	TriggerType TriggerType
	OldState    *user.User
	NewState    *user.User
}

func (a UserAction) String() string {
	return fmt.Sprintf("UserActioner for rule '%s' triggered action om User '%s' - status: %+v\n",
		a.RuleName, getUserUID(a.OldState, a.NewState), a.TriggerType)
}

func getUserUID(oldState *user.User, newState *user.User) string {
	if oldState != nil {
		return oldState.UID
	}
	return newState.UID
}
