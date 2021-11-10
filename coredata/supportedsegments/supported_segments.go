package supportedsegments

import (
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
)

const (
	OldAgeSegmentName   = "OldAgeSegment"
	YoungAgeSegmentName = "YoungAgeSegment"
)

var (
	OldAgeSegment = segment.Spec{
		UID:            OldAgeSegmentName,
		Description:    "old users segment",
		UserFilterName: user.FilterOldAgeName,
	}

	YoungAgeSegment = segment.Spec{
		UID:            YoungAgeSegmentName,
		Description:    "young users segment",
		UserFilterName: user.FilterYoungAgeName,
	}
)
