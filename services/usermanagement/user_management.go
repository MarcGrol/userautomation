package usermanagement

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/coredata/predefinedusers"
	"github.com/gorilla/mux"
	"reflect"

	"github.com/MarcGrol/userautomation/core/user"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type service struct {
	userStore     datastore.Datastore
	filterService user.FilterManager
	pubsub        pubsub.Pubsub
}

func New(store datastore.Datastore, filterService user.FilterManager, pubsub pubsub.Pubsub) user.Management {
	store.EnforceDataType(reflect.TypeOf(user.User{}).Name())
	return &service{
		userStore:     store,
		filterService: filterService,
		pubsub:        pubsub,
	}
}

func (s *service) Put(ctx context.Context, u user.User) error {
	// About publication of event:
	// - Should be published inside transaction? When if commit fails?
	// - Should be published outside transaction? When if publish fails?
	// - Maybe the event should just be stored as part of the transaction
	//     and be published by a dedicated process (and retry if publication failed)

	var eventToPublish interface{} = nil
	err := s.userStore.RunInTransaction(ctx, func(ctx context.Context) error {
		originalUser, exists, err := s.userStore.Get(ctx, u.UID)
		if err != nil {
			return fmt.Errorf("Error fetching user with uid %s: %s", u.UID, err)
		}

		err = s.userStore.Put(ctx, u.UID, u)
		if err != nil {
			return fmt.Errorf("Error storing user with uid %s: %s", u.UID, err)
		}

		if !exists {
			err := s.pubsub.Publish(ctx, user.ManagementTopicName, user.CreatedEvent{
				UserState: u,
			})
			if err != nil {
				return fmt.Errorf("Error publishing CreatedEvent for user %s: %s", u.UID, err)
			}
		} else if !reflect.DeepEqual(originalUser, u) {
			err := s.pubsub.Publish(ctx, user.ManagementTopicName, user.ModifiedEvent{
				OldUserState: originalUser.(user.User),
				NewUserState: u,
			})
			if err != nil {
				return fmt.Errorf("Error publishing ModifiedEvent for user %s: %s", u.UID, err)
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

func (s *service) Remove(ctx context.Context, userUID string) error {
	err := s.userStore.RunInTransaction(ctx, func(ctx context.Context) error {
		u, exists, err := s.userStore.Get(ctx, userUID)
		if err != nil {
			return fmt.Errorf("Error fetching user with uid %s: %s", userUID, err)
		}

		if exists {
			err = s.userStore.Remove(ctx, userUID)
			if err != nil {
				return fmt.Errorf("Error removing user with uid %s: %s", userUID, err)
			}

			err = s.pubsub.Publish(ctx, user.ManagementTopicName, user.RemovedEvent{
				UserState: u.(user.User),
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

func (s *service) Get(ctx context.Context, userUID string) (user.User, bool, error) {
	u := user.User{}
	userExists := false
	var err error

	err = s.userStore.RunInTransaction(ctx, func(ctx context.Context) error {
		found, exists, err := s.userStore.Get(ctx, userUID)
		if err != nil {
			return fmt.Errorf("Error fetching user with uid %s: %s", userUID, err)
		}
		userExists = exists

		if !exists {
			return nil
		}
		u = found.(user.User)
		return nil
	})
	if err != nil {
		return u, false, err
	}

	return u, userExists, nil
}

func (s *service) Query(ctx context.Context, filterName string) ([]user.User, error) {
	filterFunc, found := s.filterService.GetUserFilterByName(ctx, filterName)
	if !found {
		return []user.User{}, fmt.Errorf("Filter with name %s does not exist", filterName)
	}
	return s.QueryByFunc(ctx, filterFunc)
}

func (s *service) QueryByFunc(ctx context.Context, filterFunc user.FilterFunc) ([]user.User, error) {
	users := []user.User{}
	var err error

	err = s.userStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.userStore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all users: %s", err)
		}

		for _, u := range items {
			user := u.(user.User)
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

func (s *service) List(ctx context.Context) ([]user.User, error) {
	users := []user.User{}
	var err error

	err = s.userStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.userStore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all users: %s", err)
		}

		for _, u := range items {
			users = append(users, u.(user.User))
		}
		return nil
	})
	if err != nil {
		return users, err
	}

	return users, nil
}

func (m *service) Preprov(ctx context.Context) error {
	err := m.Put(ctx, predefinedusers.Marc)
	if err != nil {
		return err
	}

	err = m.Put(ctx, predefinedusers.Eva)
	if err != nil {
		return err
	}

	err = m.Put(ctx, predefinedusers.Pien)
	if err != nil {
		return err
	}

	return nil
}

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {

}
