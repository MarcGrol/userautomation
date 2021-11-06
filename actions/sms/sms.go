package sms

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/actions/actionutil"
	"github.com/MarcGrol/userautomation/actions/templating"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/users"
)

//go:generate mockgen -source=sms.go -destination=sms_mock.go -package=sms SmsSender
type SmsSender interface {
	Send(c context.Context, recipient, body string) error
}

type smsSender struct {
}

func NewSmsSender() SmsSender {
	return &smsSender{}
}

func (es *smsSender) Send(c context.Context, recipient, body string) error {
	fmt.Printf("send sms to address: '%s' with body: '%s'\n", recipient, body)
	return nil
}

func SmsAction(bodyTemplate string, smsClient SmsSender) rules.UserActionFunc {
	return func(ctx context.Context, ruleName string, userStatus rules.UserChangeStatus, oldState *users.User, newState *users.User) error {
		actionutil.LogFunc(ctx, ruleName, userStatus, oldState, newState)

		if userStatus == rules.UserRemoved {
			return nil
		}

		userPhoneNumber, ok := newState.Attributes["phonenumber"].(string)
		if !ok {
			return fmt.Errorf("User %+v has no phonenumber", newState)
		}

		body, err := templating.ApplyTemplate(ruleName+"-sms-body", bodyTemplate, newState.Attributes)
		if err != nil {
			return fmt.Errorf("Error creating sms body for newState %s:%s", newState.UID, err)
		}

		err = smsClient.Send(ctx, userPhoneNumber, body)
		if err != nil {
			return fmt.Errorf("Error sending sms for newState %s:%s", newState.UID, err)
		}

		return nil
	}
}
