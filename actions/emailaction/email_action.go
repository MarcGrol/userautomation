package emailaction

import (
	"context"
	"fmt"
	"log"

	"github.com/MarcGrol/userautomation/actions/templating"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
)

type EmailAction struct {
	subjectTemplate string
	bodyTemplate    string
	emailClient     emailsending.EmailSender
}

func NewEmailAction(subjectTemplate string, bodyTemplate string, emailClient emailsending.EmailSender) usertask.UserTaskExecutor {
	return &EmailAction{
		subjectTemplate: subjectTemplate,
		bodyTemplate:    bodyTemplate,
		emailClient:     emailClient,
	}
}

func (ea *EmailAction) Perform(ctx context.Context, a usertask.Spec) error {
	log.Printf("email-action: %s", a.String())

	userEmail, ok := a.User.Attributes["email_address"].(string)
	if !ok {
		return fmt.Errorf("User %+v has no emailaddress", a.User)
	}
	subject, err := templating.ApplyTemplate(a.ActionSpec.Name+"-email-subject", ea.subjectTemplate, a.User.Attributes)
	if err != nil {
		return fmt.Errorf("Error creating email subject for user %s:%s", a.User.UID, err)
	}

	body, err := templating.ApplyTemplate(a.ActionSpec.Name+"-email-body", ea.bodyTemplate, a.User.Attributes)
	if err != nil {
		return fmt.Errorf("Error creating email body for user %s:%s", a.User.UID, err)
	}

	err = ea.emailClient.Send(ctx, userEmail, subject, body)
	if err != nil {
		return fmt.Errorf("Error sending email for user %s:%s", a.User.UID, err)
	}

	return nil
}
