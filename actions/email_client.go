package actions

//go:generate mockgen -source=email_client.go -destination=email_client_mock.go -package=actions Emailer
type Emailer interface {
	Send(recipient, subject, body string) error
}
