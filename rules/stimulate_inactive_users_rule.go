package rules

import (
	"context"
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

func (r *stimulateInactiveUsersRule) Name() string {
	return "StimulateInactiveUsers"
}

func (r *stimulateInactiveUsersRule) ApplicableFor(event core.Event) bool {
	return event.EventName == "Timer"
}

func (r *stimulateInactiveUsersRule) DetermineAudience(c context.Context) ([]core.User, error) {
	users, err := r.userLookup.GetUserOnQuery(c, "publishCount == 0 && loginCount < 5")
	if err != nil {
		return []core.User{}, err
	}
	return users, nil
}

func (r *stimulateInactiveUsersRule) ApplyAction(c context.Context, user core.User) error {
	subject, err := utils.ApplyTemplate("stimulateInactiveUserRule-subject", "We are dissapointed", user.Payload)
	if err != nil {
		return err
	}
	body, err := utils.ApplyTemplate("stimulateInactiveUserRule-body", "Hi {{FirstName}}, do more", user.Payload)
	if err != nil {
		return err
	}

	err = r.emailer.Send(c, user.EmailAddress, subject, body)
	if err != nil {
		return err
	}

	return nil
}
