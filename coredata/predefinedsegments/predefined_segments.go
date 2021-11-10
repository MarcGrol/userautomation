package predefinedsegments

import (
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/coredata/predefinedfilters"
)

const (
	OldAgeSegmentName   = "OldAgeSegment"
	YoungAgeSegmentName = "YoungAgeSegment"
)

var (
	OldAgeSegment = segment.Spec{
		UID:            OldAgeSegmentName,
		Description:    "old users segment",
		UserFilterName: predefinedfilters.FilterOldAgeName,
	}

	YoungAgeSegment = segment.Spec{
		UID:            YoungAgeSegmentName,
		Description:    "young users segment",
		UserFilterName: predefinedfilters.FilterYoungAgeName,
	}
)
