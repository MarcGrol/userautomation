package realtimeactions

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
	"github.com/MarcGrol/userautomation/realtime/realtimeutil"
)

//go:generate mockgen -source=email.go -destination=email_mock.go -package=realtimeactions Emailer
type Emailer interface {
	Send(c context.Context, recipient, subject, body string) error
}

type emailSender struct {
}

func NewEmailSender() Emailer {
	return &emailSender{}
}

func (es *emailSender) Send(c context.Context, recipient, subject, body string) error {
	fmt.Printf("send email to address: '%s' with subject '%s' and body: '%s'\n", recipient, subject, body)
	return nil
}

func EmailerAction(subjectTemplate string, bodyTemplate string, emailClient Emailer) realtimecore.UserActionFunc {
	return func(ctx context.Context, ruleName string, userStatus realtimecore.UserChangeStatus, oldState *realtimecore.User, newState *realtimecore.User) error {
		logFunc(ctx, ruleName, userStatus, oldState, newState)

		if userStatus == realtimecore.UserRemoved {
			return nil
		}

		userEmail, ok := newState.Attributes["emailaddress"].(string)
		if !ok {
			return fmt.Errorf("User %+v has no emailaddress", newState)
		}
		subject, err := realtimeutil.ApplyTemplate(ruleName+"-email-subject", subjectTemplate, newState.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating email subject for newState %s:%s", newState.UID, err)
		}

		body, err := realtimeutil.ApplyTemplate(ruleName+"-email-body", bodyTemplate, newState.Attributes)
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
