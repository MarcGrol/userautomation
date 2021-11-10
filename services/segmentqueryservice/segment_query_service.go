package segmentmanagement

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type service struct {
	segmentStore datastore.Datastore
}

func New(datastore datastore.Datastore, pubsub pubsub.Pubsub) segment.Querier {
	return &service{
		segmentStore: datastore,
	}
}

func (s *service) GetUsersForSegment(ctx context.Context, segmentUID string) ([]user.User, error) {
	users := []user.User{}
	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.segmentStore.Get(ctx, segmentUID)
		if err != nil {
			return err
		}
		if !exists {
			return fmt.Errorf("Spec with uid %s does not exist", segmentUID)
		}

		swu := item.(segment.SegmentWithUsers)
		users := []user.User{}
		for _, u := range swu.Users {
			users = append(users, u)
		}
		return nil
	})
	if err != nil {
		return users, err
	}
	return users, nil
}
