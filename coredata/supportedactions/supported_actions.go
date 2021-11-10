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
	SmsToYoung = action.ActionSpec{
		Name:                    SmsToYoungName,
		Description:             "Sms to young people",
		MandatoryUserAttributes: []string{PhoneNumber, FirstName, Age},
		ProvidedAttributes: map[string]string{
			"body_template": "my sms body template",
		},
	}
	MailToOld = action.ActionSpec{
		Name:                    MailToOldName,
		Description:             "Mail to old people",
		MandatoryUserAttributes: []string{EmailAddress, FirstName, Age},
		ProvidedAttributes: map[string]string{
			"subject_template": "my email subject",
			"body_template":    "my email body",
		},
	}
)
