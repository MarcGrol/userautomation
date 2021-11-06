package emailaction

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/useractions/templating"
)

func EmailerAction(subjectTemplate string, bodyTemplate string, emailClient emailsending.EmailSender) rules.UserActionFunc {
	return func(ctx context.Context, action rules.UserAction) error {

		if action.UserChangeType == rules.UserRemoved {
			return nil
		}

		userEmail, ok := action.NewState.Attributes["emailaddress"].(string)
		if !ok {
			return fmt.Errorf("User %+v has no emailaddress", action.NewState)
		}
		subject, err := templating.ApplyTemplate(action.RuleName+"-email-subject", subjectTemplate, action.NewState.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating email subject for newState %s:%s", action.NewState.UID, err)
		}

		body, err := templating.ApplyTemplate(action.RuleName+"-email-body", bodyTemplate, action.NewState.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating email body for newState %s:%s", action.NewState.UID, err)
		}

		err = emailClient.Send(ctx, userEmail, subject, body)
		if err != nil {
			return fmt.Errorf("Error sending email for newState %s:%s", action.NewState.UID, err)
		}

		return nil
	}
}
