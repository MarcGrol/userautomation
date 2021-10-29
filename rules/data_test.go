package rules

import "github.com/MarcGrol/userautomation/core"

var testUser = core.User{
	UserUID:      "123",
	EmailAddress: "123@work.nl",
	PhoneNumber:  "+31612345678",
	CommunityUID: "xebia",
	Payload: map[string]interface{}{
		"FirstName": "Marc",
	},
}
