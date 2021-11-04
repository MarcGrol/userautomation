package realtimeactions

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
	"github.com/MarcGrol/userautomation/realtime/realtimeutil"
)

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
	return func(ctx context.Context, ruleName string, userStatus realtimecore.UserStatus, user realtimecore.User) error {
		logFunc(ctx, ruleName, userStatus, user)

		if userStatus == realtimecore.UserRemoved {
			return nil
		}

		userEmail, ok := user.Attributes["emailaddress"].(string)
		if !ok {
			return fmt.Errorf("User %+v has no emailaddress", user)
		}
		subject, err := realtimeutil.ApplyTemplate(ruleName+"-email-subject", subjectTemplate, user.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating email subject for user %s:%s", user.UID, err)
		}

		body, err := realtimeutil.ApplyTemplate(ruleName+"-email-body", bodyTemplate, user.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating email body for user %s:%s", user.UID, err)
		}

		err = emailClient.Send(ctx, userEmail, subject, body)
		if err != nil {
			return fmt.Errorf("Error sending email for user %s:%s", user.UID, err)
		}

		return nil
	}
}
