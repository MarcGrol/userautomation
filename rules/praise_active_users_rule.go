package rules

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/api"
	"github.com/MarcGrol/userautomation/userlookup"
)

type praiseActiveUsersRule struct {
	userLookup userlookup.UserLookuper
	emailer    actions.Emailer
}

func NewPraiseActiveUserRule(userLookup userlookup.UserLookuper, emailer actions.Emailer) api.UserRule {
	return &praiseActiveUsersRule{
		userLookup: userLookup,
		emailer:    emailer,
	}
}

func (r praiseActiveUsersRule) ApplicableFor(event api.Event) bool {
	return event.EventName == "Timer"
}

func (r praiseActiveUsersRule) DetermineAudience() ([]api.User, error) {
	users, err := r.userLookup.GetUserOnQuery("publishCount > 10 && loginCount > 20")
	if err != nil {
		return []api.User{}, err
	}
	return users, nil
}

func (r praiseActiveUsersRule) ApplyAction(users []api.User) error {
	for _, u := range users {
		subject, err := actions.ApplyTemplate("praiseActiveUsersRule-subject", "We praise your activity", u.Payload)
		if err != nil {
			return err
		}
		body, err := actions.ApplyTemplate("praiseActiveUsersRule-body", "Hi {{FirstName}}, well done", u.Payload)
		if err != nil {
			return err
		}

		err = r.emailer.Send(u.EmailAddress, subject, body)
		if err != nil {
			return err
		}
	}
	return nil
}
