package batchrules

import (
	"context"
	actions2 "github.com/MarcGrol/userautomation/batch/batchactions"
	"github.com/MarcGrol/userautomation/batch/batchcore"
	"github.com/MarcGrol/userautomation/batch/userlookup"
)

func GetUserRules(userLookup userlookup.UserLookuper, userGrouper actions2.GroupApi, emailer actions2.Emailer) []batchcore.UserRule {
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
