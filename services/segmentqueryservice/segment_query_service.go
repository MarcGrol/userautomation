package segmentmanagement

import (
	"context"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/services/segmentusermanagement"
)

type service struct {
	sum segmentusermanagement.SegmentUserManager
}

func New(sum segmentusermanagement.SegmentUserManager) segment.Querier {
	return &service{
		sum: sum,
	}
}

func (s *service) GetUsersForSegment(ctx context.Context, segmentUID string) ([]user.User, error) {
	return s.sum.GetUsersForSegment(ctx, segmentUID)
}
