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
	Name() string
	ApplicableFor(event Event) bool
	DetermineAudience() ([]User, error)
	ApplyAction(user User) error
}

func EvaluateUserRule(r UserRule, event Event) error {
	if !r.ApplicableFor(event) {
		log.Printf("Event '%s' not supported by rule '%s'", event.EventName, r.Name())
		return nil
	}

	users, err := r.DetermineAudience()
	if err != nil {
		log.Printf("Error fetching audience for rule '%s': %s", r.Name(), err)
		return err
	}

	for _, user := range users {
		err = r.ApplyAction(user)
		if err != nil {
			log.Printf("Error applying action to user '%s' for rule '%s': %s", user.UserUID, r.Name(), err)
			return err
		}
		log.Printf("Successfully applied action to user '%s' for rule '%s'", user.UserUID, r.Name())
	}

	return nil
}
