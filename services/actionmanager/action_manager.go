package actionmanager

import (
	"context"
	"github.com/MarcGrol/userautomation/core/action"
)

type actionManager struct {
	actions []action.ActionSpec
}

type ActionManager interface {
	action.ActionManager
}

func New() ActionManager {
	return &actionManager{
		actions: []action.ActionSpec{
			{
				Name:                    "SmsToYoung",
				Description:             "Sms to young people",
				MandatoryUserAttributes: []string{"phone_number", "first_name", "age"},
				ProvidedAttributes: map[string]string{
					"body_template":"my sms body template",
				},
			},
			{
				Name:                    "MailToOld",
				Description:             "Mail to old people",
				MandatoryUserAttributes: []string{"email_address", "first_name", "age"},
				ProvidedAttributes: map[string]string{
					"subject_template":"my email subject",
					"body_template":"my email body",
				},},
		},
	}
}

func(m *actionManager) GetActionSpecOnName(ctx context.Context, name string) (action.ActionSpec, bool, error) {
	for _, a := range m.actions {
		if a.Name == name {
			return a, true, nil
		}
	}

	return action.ActionSpec{}, false, nil
}

func(m *actionManager) ListActionSpecs(ctx context.Context) ([]action.ActionSpec, error) {
	return m.actions, nil
}
