package actionmanager

import (
	"context"
	"github.com/MarcGrol/userautomation/core/action"
)

type actionManager struct {
}

type ActionManager interface {
	action.ActionManager
}

func New() ActionManager {
	return &actionManager{}
}

const (
	SmsToYoungName = "SmsToYoung"
	MailToOldName  = "MailToOld"
)

var (
	SmsToYoung = action.ActionSpec{
		Name:                    SmsToYoungName,
		Description:             "Sms to young people",
		MandatoryUserAttributes: []string{"phone_number", "first_name", "age"},
		ProvidedAttributes: map[string]string{
			"body_template": "my sms body template",
		},
	}
	MailToOld = action.ActionSpec{
		Name:                    MailToOldName,
		Description:             "Mail to old people",
		MandatoryUserAttributes: []string{"email_address", "first_name", "age"},
		ProvidedAttributes: map[string]string{
			"subject_template": "my email subject",
			"body_template":    "my email body",
		},
	}
)

var actionMap = map[string]action.ActionSpec{
	SmsToYoungName: SmsToYoung,
	MailToOldName:  MailToOld,
}

func (m *actionManager) GetActionSpecOnName(ctx context.Context, name string) (action.ActionSpec, bool, error) {
	a, exists := actionMap[name]

	return a, exists, nil
}

func (m *actionManager) ListActionSpecs(ctx context.Context) ([]action.ActionSpec, error) {
	actions := []action.ActionSpec{}
	for _, a := range actionMap {
		actions = append(actions, a)
	}
	return actions, nil
}
