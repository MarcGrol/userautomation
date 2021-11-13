package userrule

import (
	"context"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/core/util"
)

type Spec struct {
	UID         string
	Description string
	User        user.User
	ActionSpec  action.Spec
}

type TriggerRuleExecution interface {
	Trigger(ctx context.Context, rule Spec) error
	util.WebExposer
}
