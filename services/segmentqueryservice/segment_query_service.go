package segmentmanagement

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type segmentQuery struct {
	segmentStore datastore.Datastore
}

type SegmentManagement interface {
	segment.Querier
}

func New(datastore datastore.Datastore, pubsub pubsub.Pubsub) SegmentManagement {
	return &segmentQuery{
		segmentStore: datastore,
	}
}

func (s *segmentQuery) GetUsersForSegment(ctx context.Context, userSegmentUID string) ([]user.User, error) {
	users := []user.User{}
	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.segmentStore.Get(ctx, userSegmentUID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Segment with uid %s does not exist", userSegmentUID)
		}

		segm := item.(segment.UserSegment)
		users := []user.User{}
		for _, u := range segm.Users {
			users = append(users, u)
		}
		return nil
	})
	if err != nil {
		return users, err
	}
	return users, nil
}