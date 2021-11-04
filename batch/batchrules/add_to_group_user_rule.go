package batchrules

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/batch/batchactions"
	"github.com/MarcGrol/userautomation/batch/batchcore"
	"github.com/MarcGrol/userautomation/batch/userlookup"
	"strings"
)

type addToGroupUserRule struct {
	userUID    string
	userLookup userlookup.UserLookuper
	groupApi   batchactions.GroupApi
}

func NewAddToGroupUserRule(userLookup userlookup.UserLookuper, groupApi batchactions.GroupApi) batchcore.UserRule {
	return &addToGroupUserRule{
		userLookup: userLookup,
		groupApi:   groupApi,
	}
}

func (r *addToGroupUserRule) Name() string {
	return "AddUserToGroup"
}

func (r *addToGroupUserRule) ApplicableFor(event batchcore.Event) bool {
	r.userUID = event.UserUID
	return event.EventName == "UserRegistered"
}

func (r *addToGroupUserRule) DetermineAudience(c context.Context) ([]batchcore.User, error) {
	user, err := r.userLookup.GetUserOnUid(c, r.userUID)
	if err != nil {
		return nil, err
	}
	return []batchcore.User{user}, nil
}

func (r *addToGroupUserRule) ApplyAction(c context.Context, user batchcore.User) error {
	groupName, err := r.deriveGroupFromEmailDomain(user.EmailAddress)
	if err != nil {
		return err
	}

	exists, err := r.groupApi.GroupExists(c, groupName)
	if err != nil {
		return err
	}

	if exists {
		err = r.groupApi.AddUserToGroup(c, groupName, user.UserUID)
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
