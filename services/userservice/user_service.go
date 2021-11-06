package userservice

import (
	"context"
	"fmt"
	user2 "github.com/MarcGrol/userautomation/core/user"
	"reflect"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type userService struct {
	datastore datastore.Datastore
	pubsub    pubsub.Pubsub
}

func NewUserService(datastore datastore.Datastore, pubsub pubsub.Pubsub) user2.Service {
	return &userService{
		datastore: datastore,
		pubsub:    pubsub,
	}
}

func (s *userService) Put(ctx context.Context, user user2.User) error {
	// About publication of event:
	// - Should be published inside transaction? When if commit fails?
	// - Should be published outside transaction? When if publish fails?
	// - Maybe the event should just be stored as part of the transaction
	//     and be published by a dedicated process (and retry if publication failed)

	var eventToPublish interface{} = nil
	err := s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		originalUser, exists, err := s.datastore.Get(ctx, user.UID)
		if err != nil {
			return fmt.Errorf("Error fetching user with uid %s: %s", user.UID, err)
		}

		err = s.datastore.Put(ctx, user.UID, user)
		if err != nil {
			return fmt.Errorf("Error storing user with uid %s: %s", user.UID, err)
		}

		if !exists {
			err := s.pubsub.Publish(ctx, user2.UserTopicName, user2.CreatedEvent{
				State: user,
			})
			if err != nil {
				return fmt.Errorf("Error publishing CreatedEvent for user %s: %s", user.UID, err)
			}
		} else if !reflect.DeepEqual(originalUser, user) {
			err := s.pubsub.Publish(ctx, user2.UserTopicName, user2.ModifiedEvent{
				OldState: originalUser.(user2.User),
				NewState: user,
			})
			if err != nil {
				return fmt.Errorf("Error publishing ModifiedEvent for user %s: %s", user.UID, err)
			}
		} else {
			// user unchanged: do not notify
		}
		return nil
	})
	if err != nil {
		return err
	}

	if eventToPublish != nil {

	}

	return nil
}

func (s *userService) Remove(ctx context.Context, userUID string) error {
	err := s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		user, exists, err := s.datastore.Get(ctx, userUID)
		if err != nil {
			return fmt.Errorf("Error fetching user with uid %s: %s", userUID, err)
		}

		if exists {
			err = s.datastore.Remove(ctx, userUID)
			if err != nil {
				return fmt.Errorf("Error removing user with uid %s: %s", userUID, err)
			}

			err = s.pubsub.Publish(ctx, user2.UserTopicName, user2.RemovedEvent{
				State: user.(user2.User),
			})
			if err != nil {
				return fmt.Errorf("Error publishing RemovedEvent for user %s: %s", userUID, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Get(ctx context.Context, userUID string) (user2.User, bool, error) {
	user := user2.User{}
	userExists := false
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		found, exists, err := s.datastore.Get(ctx, userUID)
		if err != nil {
			return fmt.Errorf("Error fetching user with uid %s: %s", userUID, err)
		}
		user = found.(user2.User)
		userExists = exists

		return nil
	})
	if err != nil {
		return user, false, err
	}

	return user, userExists, nil
}

func (s *userService) Query(ctx context.Context, filterFunc user2.FilterFunc) ([]user2.User, error) {
	users := []user2.User{}
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.datastore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all users: %s", err)
		}

		for _, u := range items {
			user := u.(user2.User)
			match, err := filterFunc(ctx, user)
			if err != nil {
				return err
			}
			if match {
				users = append(users, user)
			}
		}
		return nil
	})
	if err != nil {
		return users, err
	}

	return users, nil
}
