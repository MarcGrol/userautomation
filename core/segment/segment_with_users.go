package segment

import (
	"github.com/MarcGrol/userautomation/core/user"
)

type SegmentWithUsers struct {
	Segment UserSegment
	Users   map[string]user.User
}
