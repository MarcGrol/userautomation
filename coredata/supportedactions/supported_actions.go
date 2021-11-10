package supportedactions

import (
	"github.com/MarcGrol/userautomation/core/action"
	. "github.com/MarcGrol/userautomation/coredata/supportedattrs"
)

const (
	SmsToYoungName = "SmsToYoung"
	MailToOldName  = "MailToOld"
)

var (
	SmsToYoung = action.Spec{
		Name:                    SmsToYoungName,
		Description:             "Sms to young people",
		MandatoryUserAttributes: []string{PhoneNumber, FirstName, Age},
		ProvidedInformation: map[string]string{
			"body_template": "Message to {{.first_name}}",
		},
	}
	MailToOld = action.Spec{
		Name:                    MailToOldName,
		Description:             "Mail to old people",
		MandatoryUserAttributes: []string{EmailAddress, FirstName, Age},
		ProvidedInformation: map[string]string{
			"subject_template": "Your age is {{.age}}",
			"body_template":    "Hi {{.first_name}}, your age is {{.age}}",
		},
	}
)
