package actions

import "context"

//go:generate mockgen -source=email_client.go -destination=email_client_mock.go -package=actions Emailer
type Emailer interface {
	Send(c context.Context, recipient, subject, body string) error
}
