package rules

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/core"
	"github.com/MarcGrol/userautomation/userlookup"
	"github.com/MarcGrol/userautomation/utils"
)

type stimulateInactiveUsersRule struct {
	userLookup userlookup.UserLookuper
	emailer    actions.Emailer
}

func NewStimulateInactiveUserRule(userLookup userlookup.UserLookuper, emailer actions.Emailer) core.UserRule {
	return &stimulateInactiveUsersRule{
		userLookup: userLookup,
		emailer:    emailer,
	}
}

func (r *stimulateInactiveUsersRule) ApplicableFor(event core.Event) bool {
	return event.EventName == "Timer"
}

func (r *stimulateInactiveUsersRule) DetermineAudience() ([]core.User, error) {
	users, err := r.userLookup.GetUserOnQuery("publishCount == 0 && loginCount < 5")
	if err != nil {
		return []core.User{}, err
	}
	return users, nil
}

func (r *stimulateInactiveUsersRule) ApplyAction(users []core.User) error {
	for _, u := range users {
		subject, err := utils.ApplyTemplate("stimulateInactiveUserRule-subject", "We are dissapointed", u.Payload)
		if err != nil {
			return err
		}
		body, err := utils.ApplyTemplate("stimulateInactiveUserRule-body", "Hi {{FirstName}}, do more", u.Payload)
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
