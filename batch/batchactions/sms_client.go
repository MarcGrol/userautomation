package batchactions

import "context"

//go:generate mockgen -source=sms_client.go -destination=sms_client_mock.go -package=actions SmsSender
type SmsSender interface {
	Send(c context.Context, recipient, body string) error
}
