package smssending

import "context"

//go:generate mockgen -source=api.go -destination=sms_sender_mock.go -package=smssending SmsSender
type SmsSender interface {
	Send(c context.Context, recipient, body string) error
}
