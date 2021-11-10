package smsaction

import (
	"context"
	"fmt"
	"log"

	"github.com/MarcGrol/userautomation/actions/templating"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/integrations/smssending"
)

type SmsAction struct {
	bodyTemplate string
	smsClient    smssending.SmsSender
}

func New(bodyTemplate string, smsClient smssending.SmsSender) usertask.UserTaskExecutor {
	return &SmsAction{
		bodyTemplate: bodyTemplate,
		smsClient:    smsClient,
	}
}

func (ea *SmsAction) Perform(ctx context.Context, a usertask.Spec) (string, error) {
	log.Printf("email-action: %s", a.String())

	userPhoneNumber, ok := a.User.Attributes["phone_number"].(string)
	if !ok {
		return "", fmt.Errorf("User %+v has no phone_number", a.User)
	}

	body, err := templating.ApplyTemplate(a.ActionSpec.Name+"-sms-body", ea.bodyTemplate, a.User.Attributes)
	if err != nil {
		return "", fmt.Errorf("Error creating sms body for user %s:%s", a.User.UID, err)
	}

	err = ea.smsClient.Send(ctx, userPhoneNumber, body)
	if err != nil {
		return "", fmt.Errorf("Error sending sms for user %s:%s", a.User.UID, err)
	}

	return fmt.Sprintf("Sms with content '%s' has been sent to user '%s'", body, userPhoneNumber), nil
}
