package realtimeservices

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
	"reflect"
	"sync"
)

type userService struct {
	sync.Mutex
	users            map[string]realtimecore.User
	pubsub           realtimecore.Pubsub
}

func NewUserService(pubsub realtimecore.Pubsub) realtimecore.UserService {
	return &userService{
		users:            map[string]realtimecore.User{},
		pubsub:           pubsub,
	}
}

func (s *userService) Put(ctx context.Context, user realtimecore.User) error {
	s.Lock()
	defer s.Unlock()

	originalUser, exists, err := s.getUnlocked(ctx, user.UID)
	if err != nil {
		return fmt.Errorf("Error fetching user with uid %s: %s", user.UID, err)
	}

	s.users[user.UID] = user

	if !exists {
		err := s.pubsub.Publish(ctx, "user", realtimecore.UserCreatedEvent{
			State:user,
		})
		if err != nil {
			return fmt.Errorf("Error publishing UserCreatedEvent for user %s: %s", user.UID, err)
		}
	} else if !reflect.DeepEqual(originalUser, user) {
		err := s.pubsub.Publish(ctx, "user", realtimecore.UserModifiedEvent{
			OldState:originalUser,
			NewState: user,
		})
		if err != nil {
			return fmt.Errorf("Error publishing UserModifiedEvent for user %s: %s", user.UID, err)
		}
	} else {
		// user unchanged
	}


	return nil
}

func (s *userService) Get(ctx context.Context, userUID string) (realtimecore.User, bool, error) {
	s.Lock()
	defer s.Unlock()

	return s.getUnlocked(ctx,userUID)
}

func (s *userService) getUnlocked(ctx context.Context, userUID string) (realtimecore.User, bool, error) {
	user, exists := s.users[userUID]
	return user, exists, nil
}

func (s *userService) Query(ctx context.Context, filterFunc realtimecore.UserFilterFunc) ([]realtimecore.User, error) {
	s.Lock()
	defer s.Unlock()

	result := []realtimecore.User{}

	for _, u := range s.users {
		match, err := filterFunc(ctx, u)
		if err != nil {
			return []realtimecore.User{}, err
		}
		if match {
			result = append(result, u)
		}
	}

	return result, nil
}

func (s *userService) Delete(ctx context.Context, userUID string) error {
	s.Lock()
	defer s.Unlock()

	user, exists := s.users[userUID]
	if exists {
		delete(s.users, userUID)

		err := s.pubsub.Publish(ctx, "user", realtimecore.UserRemovedEvent{
			State:user,
		})
		if err != nil {
			return fmt.Errorf("Error publishing UserRemovedEvent for user %s: %s", user.UID, err)
		}
	}

	return nil
}
