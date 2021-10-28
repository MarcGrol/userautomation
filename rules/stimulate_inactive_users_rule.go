package rules

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/api"
	"github.com/MarcGrol/userautomation/userlookup"
)

type stimulateInactiveUsersRule struct {
	userLookup userlookup.UserLookuper
	emailer    actions.Emailer
}

func NewStimulateInactiveUserRule(userLookup userlookup.UserLookuper, emailer actions.Emailer) api.UserRule {
	return &stimulateInactiveUsersRule{
		userLookup: userLookup,
		emailer:    emailer,
	}
}

func (r stimulateInactiveUsersRule) ApplicableFor(event api.Event) bool {
	return event.EventName == "Timer"
}

func (r stimulateInactiveUsersRule) DetermineAudience() ([]api.User, error) {
	users, err := r.userLookup.GetUserOnQuery("publishCount == 0 && loginCount < 5")
	if err != nil {
		return []api.User{}, err
	}
	return users, nil
}

func (r stimulateInactiveUsersRule) ApplyAction(users []api.User) error {
	for _, u := range users {
		subject, err := actions.ApplyTemplate("stimulateInactiveUserRule-subject", "We are dissapointed", u.Payload)
		if err != nil {
			return err
		}
		body, err := actions.ApplyTemplate("stimulateInactiveUserRule-body", "Hi {{FirstName}}, do more", u.Payload)
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
