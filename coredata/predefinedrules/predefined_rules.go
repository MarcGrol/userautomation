package predefinedrules

import (
	"github.com/MarcGrol/userautomation/core/segmentrule"
	supportedrules "github.com/MarcGrol/userautomation/coredata/predefinedsegments"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
)

const (
	OldAgeEmailRuleName        = "OldAgeEmailRule"
	AnotherOldAgeEmailRuleName = "AnotherOldAgeEmailRule"
	YoungAgeSmsRuleName        = "YoungAgeSmsRuleName"
	FirstNameStartsWithMName   = "FirstNameStartsWithM"
)

var (
	OldAgeEmailRule = segmentrule.Spec{
		UID:         OldAgeEmailRuleName,
		Description: "Send email to old users",
		SegmentSpec: supportedrules.OldAgeSegment,
		ActionSpec:  supportedactions.MailToOld,
	}

	AotherOldAgeEmailRule = segmentrule.Spec{
		UID:         AnotherOldAgeEmailRuleName,
		Description: "Send another email to old users",
		SegmentSpec: supportedrules.OldAgeSegment,
		ActionSpec:  supportedactions.MailToOld,
	}

	YoungAgeSmsRule = segmentrule.Spec{
		UID:         YoungAgeSmsRuleName,
		Description: "Send sms to young users",
		SegmentSpec: supportedrules.YoungAgeSegment,
		ActionSpec:  supportedactions.SmsToYoung,
	}

	FirstNameStartsWithM = segmentrule.Spec{
		UID:         FirstNameStartsWithMName,
		Description: "Send sms to users",
		SegmentSpec: supportedrules.YoungAgeSegment,
		ActionSpec:  supportedactions.SmsToYoung,
	}
)
