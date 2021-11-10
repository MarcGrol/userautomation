package segmentusermanagement

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type SegmentUserManager interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	user.EventHandler
	segment.EventHandler
	segment.Querier
}

type segmentUserManager struct {
	segmentWithUsersStore datastore.Datastore
	userService           user.Management
	filterservice         user.UserFilterResolver
	pubsub                pubsub.Pubsub
}

func New(datastore datastore.Datastore, userService user.Management, filterservice user.UserFilterResolver, pubsub pubsub.Pubsub) *segmentUserManager {
	return &segmentUserManager{
		segmentWithUsersStore: datastore,
		userService:           userService,
		filterservice:         filterservice,
		pubsub:                pubsub,
	}
}

func (s *segmentUserManager) checkInterface(sc *segmentUserManager) SegmentUserManager {
	return sc
}

func (s *segmentUserManager) IamSubscribing() {}

func (s *segmentUserManager) Subscribe(ctx context.Context) error {
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

func (s *segmentUserManager) OnEvent(ctx context.Context, topic string, event interface{}) error {
	err := user.DispatchEvent(ctx, s, topic, event)
	if err != nil {
		return segment.DispatchManagementEvent(ctx, s, topic, event)
	}
	return nil
}

func (s *segmentUserManager) OnSegmentCreated(ctx context.Context, event segment.CreatedEvent) error {
	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {
		segm := event.SegmentState

		// add all matching users to segment
		// TODO this possibly a very large task that would lock the datastore for a long time:
		// we might want to break this up with cursors into multiple tasks
		// The segment is not ready to be used in rules untill all updates have been applied
		users, err := s.userService.Query(ctx, segm.UserFilterName)
		if err != nil {
			return err
		}

		userMap := map[string]user.User{}

		for _, u := range users {
			userMap[u.UID] = u
			err = s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{SegmentUID: segm.UID, User: u})
			if err != nil {
				// what?
			}
		}

		err = s.segmentWithUsersStore.Put(ctx, segm.UID, segment.SegmentWithUsers{
			SegmentSpec: segm,
			Users:       userMap,
		})
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *segmentUserManager) OnSegmentModified(ctx context.Context, event segment.ModifiedEvent) error {
	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {
		segm := event.NewSegmentState

		item, exist, err := s.segmentWithUsersStore.Get(ctx, segm.UID)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("SegmentWithUsers with uid %s does not exist", segm.UID)
		}
		swu := item.(segment.SegmentWithUsers)
		swu.SegmentSpec = segm

		// TODO these both are possibly a very large tasks that would lock the datastoree:
		// we might want to break this up with cursors into multiple tasks
		{
			// Remove existing users that nom longer match segment
			for _, u := range swu.Users {
				applicable, err := s.isSegmentApplicableForUser(ctx, u, segm.UserFilterName)
				if err != nil {
					return err
				}
				if !applicable {
					delete(swu.Users, u.UID)
					err = s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserRemovedFromSegmentEvent{SegmentUID: segm.UID, User: u})
					if err != nil {
						// what?
					}
				}
			}
		}

		{
			// Add matching users that were not part of segment before
			users, err := s.userService.Query(ctx, event.NewSegmentState.UserFilterName)
			if err != nil {
				return err
			}
			for _, u := range users {
				_, exists := swu.Users[u.UID]
				if !exists {
					swu.Users[u.UID] = u
					s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{SegmentUID: segm.UID, User: u})
				}
			}
		}

		err = s.segmentWithUsersStore.Put(ctx, segm.UID, swu)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *segmentUserManager) isSegmentApplicableForUser(ctx context.Context, u user.User, userFilterName string) (bool, error) {
	filterFunc, found := s.filterservice.GetUserFilterByName(ctx, userFilterName)
	if !found {
		return false, fmt.Errorf("User filter with name %s was not found", userFilterName)
	}
	matched, err := filterFunc(ctx, u)
	if err != nil {
		return false, fmt.Errorf("Error filteriing user %s: %s", u.UID, err)
	}
	return matched, nil
}

func (s *segmentUserManager) OnSegmentRemoved(ctx context.Context, event segment.RemovedEvent) error {
	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {
		segm := event.SegmentState

		item, exist, err := s.segmentWithUsersStore.Get(ctx, segm.UID)
		if err != nil {
			return err
		}
		if !exist {
			return fmt.Errorf("SegmentWithUsers with uid %s does not exist", segm.UID)
		}
		swu := item.(segment.SegmentWithUsers)

		// Remove existing users that nom longer match segment
		// TODO these both are possibly a very large task: we might want to break this up with cursors into multiple tasks
		for _, u := range swu.Users {
			err := s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserRemovedFromSegmentEvent{SegmentUID: segm.UID, User: u})
			if err != nil {
				// what?
			}
		}

		err = s.segmentWithUsersStore.Remove(ctx, event.SegmentState.UID)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *segmentUserManager) OnUserCreated(ctx context.Context, event user.CreatedEvent) error {
	u := event.UserState

	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {
		// check if user must be added to items
		items, err := s.segmentWithUsersStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, item := range items {
			swu := item.(segment.SegmentWithUsers)

			isApplicable, err := s.isSegmentApplicableForUser(ctx, u, swu.SegmentSpec.UserFilterName)
			if err != nil {
				return err
			}
			if isApplicable {
				swu.Users[u.UID] = u
				err := s.segmentWithUsersStore.Put(ctx, swu.SegmentSpec.UID, swu)
				if err != nil {
					return err
				}
				err = s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{
					SegmentUID: swu.SegmentSpec.UID,
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

func (s *segmentUserManager) OnUserModified(ctx context.Context, event user.ModifiedEvent) error {
	u := event.NewUserState

	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {

		items, err := s.segmentWithUsersStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, item := range items {
			swu := item.(segment.SegmentWithUsers)
			_, found := swu.Users[u.UID]

			isApplicable, err := s.isSegmentApplicableForUser(ctx, u, swu.SegmentSpec.UserFilterName)
			if err != nil {
				return err
			}
			if found && !isApplicable {
				delete(swu.Users, u.UID)

				err := s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserRemovedFromSegmentEvent{
					SegmentUID: swu.SegmentSpec.UID,
					User:       u,
				})
				if err != nil {
					return err
				}
			} else if isApplicable {
				swu.Users[u.UID] = u
				err := s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserAddedToSegmentEvent{
					SegmentUID: swu.SegmentSpec.UID,
					User:       u,
				})
				if err != nil {
					return err
				}
			}
			err = s.segmentWithUsersStore.Put(ctx, swu.SegmentSpec.UID, swu)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *segmentUserManager) OnUserRemoved(ctx context.Context, event user.RemovedEvent) error {
	u := event.UserState

	return s.segmentWithUsersStore.RunInTransaction(ctx, func(ctx context.Context) error {
		// check if user must be added to items
		items, err := s.segmentWithUsersStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, item := range items {
			swu := item.(segment.SegmentWithUsers)

			delete(swu.Users, u.UID)
			err = s.segmentWithUsersStore.Put(ctx, swu.SegmentSpec.UID, swu)
			if err != nil {
				return err
			}
			err := s.pubsub.Publish(ctx, segment.UserTopicName, segment.UserRemovedFromSegmentEvent{
				SegmentUID: swu.SegmentSpec.UID,
				User:       u,
			})
			if err != nil {
				return err
			}
		}

		return nil
	})
}

func (s *segmentUserManager) GetUsersForSegment(ctx context.Context, segmentUID string) ([]user.User, error) {
	item, exists, err := s.segmentWithUsersStore.Get(ctx, segmentUID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, fmt.Errorf("Spec with uid %s does not exist", segmentUID)
	}
	swu := item.(segment.SegmentWithUsers)

	users := []user.User{}
	for _, u := range swu.Users {
		users = append(users, u)
	}
	return users, nil
}
