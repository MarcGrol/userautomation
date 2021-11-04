package realtimeservices

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
	"reflect"
)

type userService struct {
	users            map[string]realtimecore.User
	userEventService realtimecore.UserEventService
}

func NewUserService(service realtimecore.UserEventService) realtimecore.UserService {
	return &userService{
		users:            map[string]realtimecore.User{},
		userEventService: service,
	}
}

func (s *userService) Put(ctx context.Context, user realtimecore.User) error {
	originalUser, exists, err := s.Get(ctx, user.UID)
	if err != nil {
		return fmt.Errorf("Error fetching user with uid %s: %s", user.UID, err)
	}

	if !exists {
		err = s.userEventService.OnUserCreated(ctx, user)
		if err != nil {
			return fmt.Errorf("OnUserCreated for user %s failed: %s", user.UID, err)
		}
	} else if !reflect.DeepEqual(originalUser, user) {
		err = s.userEventService.OnUserModified(ctx, originalUser, user)
		if err != nil {
			return fmt.Errorf("OnUserModified for user %s failed: %s", user.UID, err)
		}
	} else {
		// user unchanged
	}

	s.users[user.UID] = user

	return nil
}

func (s *userService) Get(ctx context.Context, userUID string) (realtimecore.User, bool, error) {
	user, exists := s.users[userUID]
	return user, exists, nil
}

func (s *userService) Query(ctx context.Context, filterFunc realtimecore.UserFilterFunc) ([]realtimecore.User, error) {
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
	user, exists := s.users[userUID]
	if exists {
		err := s.userEventService.OnUserRemoved(ctx, user)
		if err != nil {
			return fmt.Errorf("OnUserRemoved for user %s failed: %s", user.UID, err)
		}
		delete(s.users, userUID)
	}

	return nil
}
