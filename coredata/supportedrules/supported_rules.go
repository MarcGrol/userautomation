package supportedrules

import (
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	supportedrules "github.com/MarcGrol/userautomation/coredata/supportedsegments"
)

const (
	OldAgeEmailRuleName = "OldAgeEmailRule"
	YoungAgeSmsRuleName = "YoungAgeSmsRuleName"
)

var (
	OldAgeEmailRule = segmentrule.Spec{
		UID:         OldAgeEmailRuleName,
		Description: "Send email to old users",
		SegmentSpec: supportedrules.OldAgeSegment,
		ActionSpec:  supportedactions.MailToOld,
	}

	YoungAgeSmsRule = segmentrule.Spec{
		UID:         YoungAgeSmsRuleName,
		Description: "Send sms to young users",
		SegmentSpec: supportedrules.YoungAgeSegment,
		ActionSpec:  supportedactions.SmsToYoung,
	}
)
