package users

import (
	"context"
	"fmt"
	"reflect"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type userService struct {
	datastore datastore.Datastore
	pubsub    pubsub.Pubsub
}

func NewUserService(datastore datastore.Datastore, pubsub pubsub.Pubsub) UserService {
	return &userService{
		datastore: datastore,
		pubsub:    pubsub,
	}
}

func (s *userService) Put(ctx context.Context, user User) error {
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
			err := s.pubsub.Publish(ctx, UserTopicName, UserCreatedEvent{
				State: user,
			})
			if err != nil {
				return fmt.Errorf("Error publishing UserCreatedEvent for user %s: %s", user.UID, err)
			}
		} else if !reflect.DeepEqual(originalUser, user) {
			err := s.pubsub.Publish(ctx, UserTopicName, UserModifiedEvent{
				OldState: originalUser.(User),
				NewState: user,
			})
			if err != nil {
				return fmt.Errorf("Error publishing UserModifiedEvent for user %s: %s", user.UID, err)
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

			err = s.pubsub.Publish(ctx, UserTopicName, UserRemovedEvent{
				State: user.(User),
			})
			if err != nil {
				return fmt.Errorf("Error publishing UserRemovedEvent for user %s: %s", userUID, err)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (s *userService) Get(ctx context.Context, userUID string) (User, bool, error) {
	user := User{}
	userExists := false
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		found, exists, err := s.datastore.Get(ctx, userUID)
		if err != nil {
			return fmt.Errorf("Error fetching user with uid %s: %s", userUID, err)
		}
		user = found.(User)
		userExists = exists

		return nil
	})
	if err != nil {
		return user, false, err
	}

	return user, userExists, nil
}

func (s *userService) Query(ctx context.Context, filterFunc UserFilterFunc) ([]User, error) {
	users := []User{}
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.datastore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all users: %s", err)
		}

		for _, u := range items {
			user := u.(User)
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
