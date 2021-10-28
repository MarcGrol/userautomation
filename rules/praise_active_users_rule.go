package rules

import (
	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/core"
	"github.com/MarcGrol/userautomation/userlookup"
	"github.com/MarcGrol/userautomation/utils"
)

type praiseActiveUsersRule struct {
	userLookup userlookup.UserLookuper
	emailer    actions.Emailer
}

func NewPraiseActiveUserRule(userLookup userlookup.UserLookuper, emailer actions.Emailer) core.UserRule {
	return &praiseActiveUsersRule{
		userLookup: userLookup,
		emailer:    emailer,
	}
}

func (r *praiseActiveUsersRule) ApplicableFor(event core.Event) bool {
	return event.EventName == "Timer"
}

func (r *praiseActiveUsersRule) DetermineAudience() ([]core.User, error) {
	users, err := r.userLookup.GetUserOnQuery("publishCount > 10 && loginCount > 20")
	if err != nil {
		return []core.User{}, err
	}
	return users, nil
}

func (r *praiseActiveUsersRule) ApplyAction(users []core.User) error {
	for _, u := range users {
		subject, err := utils.ApplyTemplate("praiseActiveUsersRule-subject", "We praise your activity", u.Payload)
		if err != nil {
			return err
		}
		body, err := utils.ApplyTemplate("praiseActiveUsersRule-body", "Hi {{FirstName}}, well done", u.Payload)
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
