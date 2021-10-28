package main

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/api"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/userlookup"
)

func main() {
	var userLookup userlookup.UserLookuper
	var userGrouper actions.GroupApi
	var emailer actions.Emailer

	userRules := rules.GetUserRules(userLookup, userGrouper, emailer)

	event := api.Event{
		EventName: "Timer",
		Payload:   map[string]interface{}{},
	}

	rules.EvaluateAllUserRules(userRules, event)
}
