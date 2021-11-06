package main

import (
	"context"
	"github.com/MarcGrol/userautomation/ignore/batch/batchactions"
	"github.com/MarcGrol/userautomation/ignore/batch/batchcore"
	"github.com/MarcGrol/userautomation/ignore/batch/batchrules"
	"github.com/MarcGrol/userautomation/ignore/batch/userlookup"
)

func main() {
	var userLookup userlookup.UserLookuper
	var userGrouper batchactions.GroupApi
	var emailer batchactions.Emailer

	userRules := batchrules.GetUserRules(userLookup, userGrouper, emailer)

	event := batchcore.Event{
		EventName: "Timer",
		Payload:   map[string]interface{}{},
	}

	batchrules.EvaluateAllUserRules(context.TODO(), userRules, event)
}
