package smssending

import (
	"context"
	"fmt"
)

type smsSender struct {
}

func NewSmsSender() SmsSender {
	return &smsSender{}
}

func (es *smsSender) Send(c context.Context, recipient, body string) error {
	fmt.Printf("send sms to address: '%s' with body: '%s'\n", recipient, body)

	// TODO integrate with 3rd party product like Twilio

	return nil
}
