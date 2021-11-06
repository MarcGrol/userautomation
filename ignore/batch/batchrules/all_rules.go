package batchrules

import (
	"context"
	"github.com/MarcGrol/userautomation/ignore/batch/batchactions"
	"github.com/MarcGrol/userautomation/ignore/batch/batchcore"
	"github.com/MarcGrol/userautomation/ignore/batch/userlookup"
)

func GetUserRules(userLookup userlookup.UserLookuper, userGrouper batchactions.GroupApi, emailer batchactions.Emailer) []batchcore.UserRule {
	return []batchcore.UserRule{
		NewPraiseActiveUserRule(userLookup, emailer),
		NewStimulateInactiveUserRule(userLookup, emailer),
		NewAddToGroupUserRule(userLookup, userGrouper),
	}
}

func EvaluateAllUserRules(c context.Context, userRules []batchcore.UserRule, event batchcore.Event) error {
	for _, r := range userRules {
		err := batchcore.EvaluateUserRule(c, r, event)
		if err != nil {
			return err
		}
	}
	return nil
}
