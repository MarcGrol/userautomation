package predefinedusers

import (
	"github.com/MarcGrol/userautomation/core/user"
	. "github.com/MarcGrol/userautomation/coredata/supportedattrs"
)

var (
	Marc = user.User{
		UID: "1",
		Attributes: map[string]interface{}{
			FirstName:    "Marc",
			EmailAddress: "marc@home.nl",
			PhoneNumber:  "+31611111111",
			Age:          50,
		},
	}
	Eva = user.User{
		UID: "2",
		Attributes: map[string]interface{}{
			FirstName:    "Eva",
			EmailAddress: "eva@home.nl",
			PhoneNumber:  "+31622222222",
			Age:          48,
		},
	}
	Pien = user.User{
		UID: "3",
		Attributes: map[string]interface{}{
			FirstName:    "Pien",
			EmailAddress: "pien@home.nl",
			PhoneNumber:  "+316333333",
			Age:          16,
		},
	}
)
