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
}

func (s *segmentService) OnUserModified(ctx context.Context, _ user.User, u user.User) error {
	// check if user must be added to segments
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
}

func (s *segmentService) OnUserRemoved(ctx context.Context, u user.User) error {
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
}

func (s *segmentService) Put(ctx context.Context, segm segment.UserSegment) error {
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
}

func (s *segmentService) Get(ctx context.Context, userSegmentUID string) (segment.UserSegment, bool, error) {
	item, exists, err := s.segmentStore.Get(ctx, userSegmentUID)
	return item.(segment.UserSegment), exists, err
}

func (s *segmentService) List(ctx context.Context) ([]segment.UserSegment, error) {
	items, err := s.segmentStore.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	segments := []segment.UserSegment{}
	for _, i := range items {
		segments = append(segments, i.(segment.UserSegment))
	}

	return segments, nil
}

func (s *segmentService) Remove(ctx context.Context, userSegmentUID string) error {
	// Can we remove a segment? I might still be in use by a rule

	return fmt.Errorf("Not implemented")
}

func (s *segmentService) GetUsersForSegment(ctx context.Context, userSegmentUID string) ([]user.User, error) {
	item, exists, err := s.segmentStore.Get(ctx, userSegmentUID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Segment with uid %s does not exist", userSegmentUID)
	}

	segm := item.(segment.UserSegment)
	users := []user.User{}
	for _, u := range segm.Users {
		users = append(users, u)
	}

	return users, nil
}
