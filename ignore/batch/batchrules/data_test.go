package batchrules

import (
	"github.com/MarcGrol/userautomation/ignore/batch/batchcore"
)

var testUser = batchcore.User{
	UserUID:      "123",
	EmailAddress: "123@work.nl",
	PhoneNumber:  "+31612345678",
	CommunityUID: "xebia",
	Payload: map[string]interface{}{
		"FirstName": "Marc",
	},
}
