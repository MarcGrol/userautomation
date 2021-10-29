package core

import "log"

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
	ApplyAction(user User) error
}

func EvaluateUserRule(r UserRule, event Event) error {
	if !r.ApplicableFor(event) {
		log.Printf("Event  %s not supported by rule", event.EventName)
		return nil
	}

	users, err := r.DetermineAudience()
	if err != nil {
		log.Printf("Error fetching audience for rule: %s", err)
		return err
	}

	for _, user := range users {
		err = r.ApplyAction(user)
		if err != nil {
			log.Printf("Error applying action for rule: %s", err)
			return err
		}
	}

	return nil
}
