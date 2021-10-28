package rules

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/core"
	"github.com/MarcGrol/userautomation/userlookup"
)

func GetUserRules(userLookup userlookup.UserLookuper, userGrouper actions.GroupApi, emailer actions.Emailer) []core.UserRule {
	return []core.UserRule{
		NewPraiseActiveUserRule(userLookup, emailer),
		NewStimulateInactiveUserRule(userLookup, emailer),
		NewAddToGroupUserRule(userLookup, userGrouper),
	}
}

func EvaluateAllUserRules(userRules []core.UserRule, event core.Event) error {
	for _, r := range userRules {
		err := core.EvaluateUserRule(r, event)
		if err != nil {
			return err
		}
	}
	return nil
}
