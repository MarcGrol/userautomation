package main

import (
	"context"
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/core"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/userlookup"
)

func main() {
	var userLookup userlookup.UserLookuper
	var userGrouper actions.GroupApi
	var emailer actions.Emailer

	userRules := rules.GetUserRules(userLookup, userGrouper, emailer)

	event := core.Event{
		EventName: "Timer",
		Payload:   map[string]interface{}{},
	}

	rules.EvaluateAllUserRules(context.TODO(), userRules, event)
}
