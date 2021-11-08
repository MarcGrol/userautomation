package segmentcalculator

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"log"
)

type segmentCalculator struct {
	segmentWithUsersStore datastore.Datastore
	userService           user.Management
	pubsub                pubsub.Pubsub
}

type SegmentCalculator interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	user.EventHandler
	segment.EventHandler
	segment.Querier
}

func New(datastore datastore.Datastore, userService user.Management, pubsub pubsub.Pubsub) *segmentCalculator {
	return &segmentCalculator{
		segmentWithUsersStore: datastore,
		userService:           userService,
		pubsub:                pubsub,
	}
}

func (s *segmentCalculator) checkInterface(sc *segmentCalculator) SegmentCalculator {
	return sc
}

func (s *segmentCalculator) IamSubscribing() {}

func (s *segmentCalculator) Subscribe(ctx context.Context) error {
	err := s.pubsub.Subscribe(ctx, user.ManagementTopicName, s.OnEvent)
	if err != nil {
		return err
	}
	err = s.pubsub.Subscribe(ctx, segment.ManagementTopicName, s.OnEvent)
	if err != nil {
		return err
	}
	return nil
}

func (s *segmentCalculator) OnEvent(ctx context.Context, topic string, event interface{}) error {
	err := user.DispatchEvent(ctx, s, topic, event)
	if err != nil {
		return segment.DispatchManagementEvent(ctx, s, topic, event)
	}
	return nil
}

func (s *segmentCalculator) OnSegmentCreated(ctx context.Context, event segment.CreatedEvent) error {
	segm := event.SegmentState
	users, err := s.userService.QueryByName(ctx, segm.UserFilterName)
	if err != nil {
		return err
	}
	segm.Users = map[string]user.User{}

	for _, u := range users {
		segm.Users[u.UID] = u
		s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{SegmentUID: segm.UID, User: u})
	}

	err = s.segmentWithUsersStore.Put(ctx, segm.UID, segm)
	if err != nil {
		return err
	}

	return nil
}

func (s *segmentCalculator) OnSegmentModified(ctx context.Context, event segment.ModifiedEvent) error {
	segm := event.NewSegmentState

	// Add or remove existing users of segment
	for _, u := range segm.Users {
		applicable, err := segm.IsApplicableForUser(ctx, u)
		if err != nil {
			return err
		}
		if applicable {
			segm.Users[u.UID] = u
			s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{SegmentUID: segm.UID, User: u})
		} else {
			delete(segm.Users, u.UID)
			s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserRemovedFromSegmentEvent{SegmentUID: segm.UID, User: u})
		}
	}

	// Add matching users that are still not part
	users, err := s.userService.QueryByName(ctx, event.NewSegmentState.UserFilterName)
	if err != nil {
		return err
	}
	log.Printf("Found %d matchings users in total set", len(users))
	for _, u := range users {
		_, exists := segm.Users[u.UID]
		if !exists {
			log.Printf("Found user %+v -> %+v", u, segm)
			segm.Users[u.UID] = u
			s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{SegmentUID: segm.UID, User: u})
		}
	}

	err = s.segmentWithUsersStore.Put(ctx, segm.UID, segm)
	if err != nil {
		return err
	}

	return nil
}

func (s *segmentCalculator) OnSegmentRemoved(ctx context.Context, event segment.RemovedEvent) error {
	err := s.segmentWithUsersStore.Remove(ctx, event.SegmentState.UID)
	if err != nil {
		return err
	}

	return nil
}

func (s *segmentCalculator) OnUserCreated(ctx context.Context, event user.CreatedEvent) error {
	u := event.UserState

	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {
		// check if user must be added to segments
		segments, err := s.segmentWithUsersStore.GetAll(ctx)
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
				err := s.segmentWithUsersStore.Put(ctx, segm.UID, segm)
				if err != nil {
					return err
				}
				err = s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{
					SegmentUID: segm.UID,
					User:       u,
				})
				if err != nil {
					return err
				}
			}
		}
		return nil
	})
}

func (s *segmentCalculator) OnUserModified(ctx context.Context, event user.ModifiedEvent) error {
	u := event.NewUserState

	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {

		segments, err := s.segmentWithUsersStore.GetAll(ctx)
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
				delete(segm.Users, u.UID)

				err := s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserRemovedFromSegmentEvent{
					SegmentUID: segm.UID,
					User:       u,
				})
				if err != nil {
					return err
				}
			} else if isApplicable {
				segm.Users[u.UID] = u
				err := s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{
					SegmentUID: segm.UID,
					User:       u,
				})
				if err != nil {
					return err
				}
			}
			err = s.segmentWithUsersStore.Put(ctx, segm.UID, segm)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *segmentCalculator) OnUserRemoved(ctx context.Context, event user.RemovedEvent) error {
	u := event.UserState

	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {
		// check if user must be added to segments
		segments, err := s.segmentWithUsersStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, item := range segments {
			segm := item.(segment.UserSegment)
			delete(segm.Users, u.UID)
			err = s.segmentWithUsersStore.Put(ctx, segm.UID, segm)
			if err != nil {
				return err
			}
			err := s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserRemovedFromSegmentEvent{
				SegmentUID: segm.UID,
				User:       u,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *segmentCalculator) GetUsersForSegment(ctx context.Context, segmentUID string) ([]user.User, error) {
	item, exists, err := s.segmentWithUsersStore.Get(ctx, segmentUID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Segment with uid %s does not exist", segmentUID)
	}
	segm := item.(segment.UserSegment)

	users := []user.User{}
	for _, u := range segm.Users {
		users = append(users, u)
	}
	return users, nil
}
