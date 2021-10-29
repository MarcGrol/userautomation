package rules

import (
	"fmt"
	"strings"

	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/core"
	"github.com/MarcGrol/userautomation/userlookup"
)

type addToGroupUserRule struct {
	userUID    string
	userLookup userlookup.UserLookuper
	groupApi   actions.GroupApi
}

func NewAddToGroupUserRule(userLookup userlookup.UserLookuper, groupApi actions.GroupApi) core.UserRule {
	return &addToGroupUserRule{
		userLookup: userLookup,
		groupApi:   groupApi,
	}
}

func (r *addToGroupUserRule) Name() string {
	return "AddUserToGroup"
}

func (r *addToGroupUserRule) ApplicableFor(event core.Event) bool {
	r.userUID = event.UserUID
	return event.EventName == "UserRegistered"
}

func (r *addToGroupUserRule) DetermineAudience() ([]core.User, error) {
	user, err := r.userLookup.GetUserOnUid(r.userUID)
	if err != nil {
		return nil, err
	}
	return []core.User{user}, nil
}

func (r *addToGroupUserRule) ApplyAction(user core.User) error {
	groupName, err := r.deriveGroupFromEmailDomain(user.EmailAddress)
	if err != nil {
		return err
	}

	exists, err := r.groupApi.GroupExists(groupName)
	if err != nil {
		return err
	}

	if exists {
		err = r.groupApi.AddUserToGroup(groupName, user.UserUID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *addToGroupUserRule) deriveGroupFromEmailDomain(userEmail string) (string, error) {
	parts := strings.Split(userEmail, "@")
	if len(parts) != 2 {
		return "", fmt.Errorf("Email address %s cannot be splitted", userEmail)
	}
	return parts[1], nil
}
