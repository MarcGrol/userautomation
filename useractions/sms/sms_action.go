package sms

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/integrations/smssending"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/useractions/templating"
)

func SmsAction(bodyTemplate string, smsClient smssending.SmsSender) rules.UserActionFunc {
	return func(ctx context.Context, action rules.UserAction) error {

		if action.UserChangeType == rules.UserRemoved {
			return nil
		}

		userPhoneNumber, ok := action.NewState.Attributes["phonenumber"].(string)
		if !ok {
			return fmt.Errorf("User %+v has no phonenumber", action.NewState)
		}

		body, err := templating.ApplyTemplate(action.RuleName+"-sms-body", bodyTemplate, action.NewState.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating sms body for newState %s:%s", action.NewState.UID, err)
		}

		err = smsClient.Send(ctx, userPhoneNumber, body)
		if err != nil {
			return fmt.Errorf("Error sending sms for newState %s:%s", action.NewState.UID, err)
		}

		return nil
	}
}
