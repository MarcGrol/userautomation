package rules

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/api"
	"github.com/MarcGrol/userautomation/userlookup"
)

type addToGroupUserRule struct {
	userUID    string
	userLookup userlookup.UserLookuper
	groupApi   actions.GroupApi
}

func NewAddToGroupUserRule(userLookup userlookup.UserLookuper, groupApi actions.GroupApi) api.UserRule {
	return &addToGroupUserRule{
		userLookup: userLookup,
		groupApi:   groupApi,
	}
}

func (r addToGroupUserRule) ApplicableFor(event api.Event) bool {
	r.userUID = event.UserUID
	return event.EventName == "UserRegistered"
}

func (r addToGroupUserRule) DetermineAudience() ([]api.User, error) {
	user, err := r.userLookup.GetUserOnUid(r.userUID)
	if err != nil {
		return nil, err
	}
	return []api.User{user}, nil
}

func (r addToGroupUserRule) ApplyAction(users []api.User) error {
	for _, u := range users {
		groupName := r.deriveGroupFromEmailDomain(u)

		exists, err := r.groupApi.GroupExists(groupName)
		if err != nil {
			return err
		}
		if exists {
			err = r.groupApi.AddUserToGroup(groupName, u.UserUID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (r addToGroupUserRule) deriveGroupFromEmailDomain(user api.User) string {
	return "tesla.com"
}
