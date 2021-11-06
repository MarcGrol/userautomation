package batchrules

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/ignore/batch/batchactions"
	"github.com/MarcGrol/userautomation/ignore/batch/batchcore"
	"github.com/MarcGrol/userautomation/ignore/batch/userlookup"
	"strings"
)

type addToGroupUserRule struct {
	userUID    string
	userLookup userlookup.UserLookuper
	groupAPI   batchactions.GroupAPI
}

func NewAddToGroupUserRule(userLookup userlookup.UserLookuper, groupAPI batchactions.GroupAPI) batchcore.UserRule {
	return &addToGroupUserRule{
		userLookup: userLookup,
		groupAPI:   groupAPI,
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

	exists, err := r.groupAPI.GroupExists(c, groupName)
	if err != nil {
		return err
	}

	if exists {
		err = r.groupAPI.AddUserToGroup(c, groupName, user.UserUID)
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
