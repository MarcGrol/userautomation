package batchrules

import (
	"context"
	"github.com/MarcGrol/userautomation/ignore/batch/batchactions"
	"github.com/MarcGrol/userautomation/ignore/batch/batchcore"
	"github.com/MarcGrol/userautomation/ignore/batch/batchutils"
	"github.com/MarcGrol/userautomation/ignore/batch/userlookup"
)

type praiseActiveUsersRule struct {
	userLookup userlookup.UserLookuper
	emailer    batchactions.Emailer
}

func NewPraiseActiveUserRule(userLookup userlookup.UserLookuper, emailer batchactions.Emailer) batchcore.UserRule {
	return &praiseActiveUsersRule{
		userLookup: userLookup,
		emailer:    emailer,
	}
}

func (r *praiseActiveUsersRule) Name() string {
	return "PraiseActiveUsers"
}

func (r *praiseActiveUsersRule) ApplicableFor(event batchcore.Event) bool {
	return event.EventName == "Timer"
}

func (r *praiseActiveUsersRule) DetermineAudience(c context.Context) ([]batchcore.User, error) {
	users, err := r.userLookup.GetUserOnQuery(c, "publishCount > 10 && loginCount > 20")
	if err != nil {
		return []batchcore.User{}, err
	}
	return users, nil
}

func (r *praiseActiveUsersRule) ApplyAction(c context.Context, user batchcore.User) error {
	subject, err := batchutils.ApplyTemplate("praiseActiveUsersRule-subject", "We praise your activity", user.Payload)
	if err != nil {
		return err
	}

	body, err := batchutils.ApplyTemplate("praiseActiveUsersRule-body", "Hi {{.FirstName}}, well done", user.Payload)
	if err != nil {
		return err
	}

	err = r.emailer.Send(c, user.EmailAddress, subject, body)
	if err != nil {
		return err
	}

	return nil
}
