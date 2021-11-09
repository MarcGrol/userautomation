package emailaction

import (
	"context"
	"fmt"
	"log"

	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/actions/templating"
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

func (ea *EmailAction) Perform(ctx context.Context, a usertask.UserTask) error {
	log.Printf("email-action: %s", a.String())

	userEmail, ok := a.NewState.Attributes["emailaddress"].(string)
	if !ok {
		return fmt.Errorf("User %+v has no emailaddress", a.NewState)
	}
	subject, err := templating.ApplyTemplate(a.RuleUID+"-email-subject", ea.subjectTemplate, a.NewState.Attributes)
	if err != nil {
		return fmt.Errorf("Error creating email subject for newState %s:%s", a.NewState.UID, err)
	}

	body, err := templating.ApplyTemplate(a.RuleUID+"-email-body", ea.bodyTemplate, a.NewState.Attributes)
	if err != nil {
		return fmt.Errorf("Error creating email body for newState %s:%s", a.NewState.UID, err)
	}

	err = ea.emailClient.Send(ctx, userEmail, subject, body)
	if err != nil {
		return fmt.Errorf("Error sending email for newState %s:%s", a.NewState.UID, err)
	}

	return nil
}
