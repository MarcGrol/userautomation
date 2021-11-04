package realtimeactions

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
	"github.com/MarcGrol/userautomation/realtime/realtimeutil"
)

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

func SmsAction(bodyTemplate string, smsClient SmsSender) realtimecore.UserActionFunc {
	return func(ctx context.Context, ruleName string, userStatus realtimecore.UserChangeStatus, oldState *realtimecore.User, newState *realtimecore.User) error {
		logFunc(ctx, ruleName, userStatus, oldState, newState)

		if userStatus == realtimecore.UserRemoved {
			return nil
		}

		userPhoneNumber, ok := newState.Attributes["phonenumber"].(string)
		if !ok {
			return fmt.Errorf("User %+v has no phonenumber", newState)
		}

		body, err := realtimeutil.ApplyTemplate(ruleName+"-sms-body", bodyTemplate, newState.Attributes)
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
