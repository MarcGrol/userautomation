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

func (r *praiseActiveUsersRule) Name() string {
	return "PraiseActiveUsers"
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

func (r *praiseActiveUsersRule) ApplyAction(user core.User) error {
	subject, err := utils.ApplyTemplate("praiseActiveUsersRule-subject", "We praise your activity", user.Payload)
	if err != nil {
		return err
	}

	body, err := utils.ApplyTemplate("praiseActiveUsersRule-body", "Hi {{.FirstName}}, well done", user.Payload)
	if err != nil {
		return err
	}

	err = r.emailer.Send(user.EmailAddress, subject, body)
	if err != nil {
		return err
	}

	return nil
}
