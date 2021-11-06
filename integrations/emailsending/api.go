package emailsending

import "context"

//go:generate mockgen -source=api.go -destination=email_sender_mock.go -package=emailsending EmailSender
type EmailSender interface {
	Send(c context.Context, recipient, subject, body string) error
}
