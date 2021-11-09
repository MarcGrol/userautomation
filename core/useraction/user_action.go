package useraction

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/user"
)

//go:generate mockgen -source=action.go -destination=actioner_mock.go -package=action UserActioner
type UserActioner interface {
	Perform(ctx context.Context, action UserAction) error
}

type ReasonForAction int

const (
	ReasonIsUserAddedToSegment ReasonForAction = iota
	ReasonIsOnDemand
	ReasonIsCron
)

type UserAction struct {
	RuleUID  string
	Reason   ReasonForAction
	OldState *user.User
	NewState *user.User
}

func (a UserAction) String() string {
	return fmt.Sprintf("UserActioner for rule '%s' triggered action om User '%s' - status: %+v\n",
		a.RuleUID, getUserUID(a.OldState, a.NewState), a.Reason)
}

func getUserUID(oldState *user.User, newState *user.User) string {
	if oldState != nil {
		return oldState.UID
	}
	return newState.UID
}
