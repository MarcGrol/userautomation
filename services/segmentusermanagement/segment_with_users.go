package segmentusermanagement

import (
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
)

type segmentWithUsers struct {
	SegmentSpec segment.Spec
	Users       map[string]user.User
}
