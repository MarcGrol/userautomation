package email

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/actions/actionutil"
	"github.com/MarcGrol/userautomation/actions/templating"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/users"
)

//go:generate mockgen -source=email.go -destination=email_mock.go -package=email EmailSender
type EmailSender interface {
	Send(c context.Context, recipient, subject, body string) error
}

type emailSender struct {
}

func NewEmailSender() EmailSender {
	return &emailSender{}
}

func (es *emailSender) Send(c context.Context, recipient, subject, body string) error {
	fmt.Printf("send email to address: '%s' with subject '%s' and body: '%s'\n", recipient, subject, body)
	return nil
}

func EmailerAction(subjectTemplate string, bodyTemplate string, emailClient EmailSender) rules.UserActionFunc {
	return func(ctx context.Context, ruleName string, userStatus rules.UserChangeStatus, oldState *users.User, newState *users.User) error {
		actionutil.LogFunc(ctx, ruleName, userStatus, oldState, newState)

		if userStatus == rules.UserRemoved {
			return nil
		}

		userEmail, ok := newState.Attributes["emailaddress"].(string)
		if !ok {
			return fmt.Errorf("User %+v has no emailaddress", newState)
		}
		subject, err := templating.ApplyTemplate(ruleName+"-email-subject", subjectTemplate, newState.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating email subject for newState %s:%s", newState.UID, err)
		}

		body, err := templating.ApplyTemplate(ruleName+"-email-body", bodyTemplate, newState.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating email body for newState %s:%s", newState.UID, err)
		}

		err = emailClient.Send(ctx, userEmail, subject, body)
		if err != nil {
			return fmt.Errorf("Error sending email for newState %s:%s", newState.UID, err)
		}

		return nil
	}
}
