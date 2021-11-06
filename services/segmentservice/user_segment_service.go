package segmentservice

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type segmentService struct {
	segmentStore datastore.Datastore
	userService  user.Service
	pubsub       pubsub.Pubsub
}

type SegmentService interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	user.EventHandler
	segment.UserSegmentService
	segment.UserSegmentQueryService
}

func NewSegmentService(datastore datastore.Datastore, userService user.Service, pubsub pubsub.Pubsub) SegmentService {
	return &segmentService{
		segmentStore: datastore,
		userService:  userService,
		pubsub:       pubsub,
	}
}

func (s *segmentService) IamSubscribing() {}

func (s *segmentService) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, user.UserTopicName, s.OnEvent)
}

func (s *segmentService) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return user.DispatchEvent(ctx, s, topic, event)
}

func (s *segmentService) OnUserCreated(ctx context.Context, u user.User) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		// check if user must be added to segments
		segments, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, item := range segments {
			segm := item.(segment.UserSegment)
			isApplicable, err := segm.IsApplicableForUser(ctx, u)
			if err != nil {
				return err
			}
			if isApplicable {
				segm.Users[u.UID] = u
				err := s.segmentStore.Put(ctx, segm.UID, segm)
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *segmentService) OnUserModified(ctx context.Context, _ user.User, u user.User) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {

		segments, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, item := range segments {
			segm := item.(segment.UserSegment)
			_, found := segm.Users[u.UID]
			isApplicable, err := segm.IsApplicableForUser(ctx, u)
			if err != nil {
				return err
			}
			if found && !isApplicable {
				err = s.segmentStore.Remove(ctx, segm.UID)
				if err != nil {
					return err
				}
			} else if isApplicable {
				segm.Users[u.UID] = u
				err = s.segmentStore.Put(ctx, segm.UID, segm)
				if err != nil {
					return err
				}

			}
		}
		return nil
	})
}

func (s *segmentService) OnUserRemoved(ctx context.Context, u user.User) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		// check if user must be added to segments
		segments, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, item := range segments {
			segm := item.(segment.UserSegment)
			delete(segm.Users, u.UID)
			err = s.segmentStore.Put(ctx, segm.UID, segm)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *segmentService) Put(ctx context.Context, segm segment.UserSegment) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		users, err := s.userService.Query(ctx, segm.IsApplicableForUser)
		if err != nil {
			return err
		}

		// recalculate list of users
		segm.Users = map[string]user.User{}
		for _, u := range users {
			applicable, err := segm.IsApplicableForUser(ctx, u)
			if err != nil {
				return err
			}
			if applicable {
				segm.Users[u.UID] = u
			}
		}

		if len(segm.Users) > 0 {
			err = s.segmentStore.Put(ctx, segm.UID, segm)
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *segmentService) Get(ctx context.Context, userSegmentUID string) (segment.UserSegment, bool, error) {
	var segm segment.UserSegment
	segmentExists := false
	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.segmentStore.Get(ctx, userSegmentUID)
		if err != nil {
			return err
		}
		segm = item.(segment.UserSegment)
		segmentExists = exists

		return nil
	})
	if err != nil {
		return segm, false, err
	}
	return segm, segmentExists, nil
}

func (s *segmentService) List(ctx context.Context) ([]segment.UserSegment, error) {
	segments := []segment.UserSegment{}

	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return err
		}

		segments := []segment.UserSegment{}
		for _, i := range items {
			segments = append(segments, i.(segment.UserSegment))
		}

		return nil
	})
	if err != nil {
		return segments, err
	}
	return segments, nil
}

func (s *segmentService) Remove(ctx context.Context, userSegmentUID string) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		// Can we remove a segment? I might still be in use by a rule

		return fmt.Errorf("Not implemented")
	})
}

func (s *segmentService) GetUsersForSegment(ctx context.Context, userSegmentUID string) ([]user.User, error) {
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
