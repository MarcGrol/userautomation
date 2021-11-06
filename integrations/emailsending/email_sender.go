package emailsending

import (
	"context"
	"fmt"
)

type emailSender struct {
}

func NewEmailSender() EmailSender {
	return &emailSender{}
}

func (es *emailSender) Send(c context.Context, recipient, subject, body string) error {
	fmt.Printf("send email to address: '%s' with subject '%s' and body: '%s'\n", recipient, subject, body)
	return nil
}
