package rules

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/api"
	"github.com/MarcGrol/userautomation/userlookup"
)

func GetUserRules(userLookup userlookup.UserLookuper, userGrouper actions.GroupApi, emailer actions.Emailer) []api.UserRule {
	return []api.UserRule{
		NewPraiseActiveUserRule(userLookup, emailer),
		NewStimulateInactiveUserRule(userLookup, emailer),
		NewAddToGroupUserRule(userLookup, userGrouper),
	}
}

func EvaluateAllUserRules(userRules []api.UserRule, event api.Event) error {
	for _, r := range userRules {
		err := api.EvaluateUserRule(r, event)
		if err != nil {
			return err
		}
	}
	return nil
}
