package core

type Event struct {
	EventName    string
	UserUID      string
	CommunityUID string
	Payload      map[string]interface{}
}

type User struct {
	UserUID      string
	EmailAddress string
	PhoneNumber  string
	CommunityUID string
	Payload      map[string]interface{}
}

type UserRule interface {
	ApplicableFor(event Event) bool
	DetermineAudience() ([]User, error)
	ApplyAction(aud []User) error
}

func EvaluateUserRule(r UserRule, event Event) error {
	if !r.ApplicableFor(event) {
		return nil
	}

	aud, err := r.DetermineAudience()
	if err != nil {
		return err
	}

	err = r.ApplyAction(aud)
	if err != nil {
		return err
	}

	return nil
}
