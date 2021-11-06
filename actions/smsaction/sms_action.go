package smsaction

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/action"

	"github.com/MarcGrol/userautomation/actions/templating"
	"github.com/MarcGrol/userautomation/integrations/smssending"
)

type SmsAction struct {
	bodyTemplate string
	smsClient    smssending.SmsSender
}

func New(bodyTemplate string, smsClient smssending.SmsSender) action.UserActioner {
	return &SmsAction{
		bodyTemplate: bodyTemplate,
		smsClient:    smsClient,
	}
}

func (ea *SmsAction) Perform(ctx context.Context, a action.UserAction) error {

	if a.UserChangeType == action.UserRemoved {
		return nil
	}

	userPhoneNumber, ok := a.NewState.Attributes["phonenumber"].(string)
	if !ok {
		return fmt.Errorf("User %+v has no phonenumber", a.NewState)
	}

	body, err := templating.ApplyTemplate(a.RuleName+"-sms-body", ea.bodyTemplate, a.NewState.Attributes)
	if err != nil {
		return fmt.Errorf("Error creating sms body for newState %s:%s", a.NewState.UID, err)
	}

	err = ea.smsClient.Send(ctx, userPhoneNumber, body)
	if err != nil {
		return fmt.Errorf("Error sending sms for newState %s:%s", a.NewState.UID, err)
	}

	return nil
}
