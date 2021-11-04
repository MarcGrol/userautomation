package main

import (
	"context"
	actions2 "github.com/MarcGrol/userautomation/batch/batchactions"
	"github.com/MarcGrol/userautomation/batch/batchcore"
	"github.com/MarcGrol/userautomation/batch/batchrules"
	"github.com/MarcGrol/userautomation/batch/userlookup"
)

func main() {
	var userLookup userlookup.UserLookuper
	var userGrouper actions2.GroupApi
	var emailer actions2.Emailer

	userRules := batchrules.GetUserRules(userLookup, userGrouper, emailer)

	event := batchcore.Event{
		EventName: "Timer",
		Payload:   map[string]interface{}{},
	}

	batchrules.EvaluateAllUserRules(context.TODO(), userRules, event)
}
