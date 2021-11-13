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

//go:generate mockgen -source=user_rule.go -destination=rule_execution_mock.go -package=userrule TriggerRuleExecution
type TriggerRuleExecution interface {
	Trigger(ctx context.Context, rule Spec) error
	util.WebExposer
}
