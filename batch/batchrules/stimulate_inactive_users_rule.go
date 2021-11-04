package batchrules

import (
	"context"
	"github.com/MarcGrol/userautomation/batch/batchactions"
	"github.com/MarcGrol/userautomation/batch/batchcore"
	"github.com/MarcGrol/userautomation/batch/batchutils"
	"github.com/MarcGrol/userautomation/batch/userlookup"
)

type stimulateInactiveUsersRule struct {
	userLookup userlookup.UserLookuper
	emailer    batchactions.Emailer
}

func NewStimulateInactiveUserRule(userLookup userlookup.UserLookuper, emailer batchactions.Emailer) batchcore.UserRule {
	return &stimulateInactiveUsersRule{
		userLookup: userLookup,
		emailer:    emailer,
	}
}

func (r *stimulateInactiveUsersRule) Name() string {
	return "StimulateInactiveUsers"
}

func (r *stimulateInactiveUsersRule) ApplicableFor(event batchcore.Event) bool {
	return event.EventName == "Timer"
}

func (r *stimulateInactiveUsersRule) DetermineAudience(c context.Context) ([]batchcore.User, error) {
	users, err := r.userLookup.GetUserOnQuery(c, "publishCount == 0 && loginCount < 5")
	if err != nil {
		return []batchcore.User{}, err
	}
	return users, nil
}

func (r *stimulateInactiveUsersRule) ApplyAction(c context.Context, user batchcore.User) error {
	subject, err := batchutils.ApplyTemplate("stimulateInactiveUserRule-subject", "We are dissapointed", user.Payload)
	if err != nil {
		return err
	}
	body, err := batchutils.ApplyTemplate("stimulateInactiveUserRule-body", "Hi {{FirstName}}, do more", user.Payload)
	if err != nil {
		return err
	}

	err = r.emailer.Send(c, user.EmailAddress, subject, body)
	if err != nil {
		return err
	}

	return nil
}
