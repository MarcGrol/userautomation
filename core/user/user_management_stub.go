package user

import (
	"context"
	"fmt"
)

type UserManagementStub struct {
	filterService UserFilterResolver
	Users         map[string]User
}

func NewUserManagementStub(filterService UserFilterResolver) *UserManagementStub {
	return &UserManagementStub{
		filterService: filterService,
		Users:         map[string]User{},
	}
}

func (s *UserManagementStub) Put(ctx context.Context, u User) error {
	s.Users[u.UID] = u
	return nil
}

func (s *UserManagementStub) Remove(ctx context.Context, uid string) error {
	delete(s.Users, uid)
	return nil
}

func (s *UserManagementStub) Get(ctx context.Context, uid string) (User, bool, error) {
	item, exists := s.Users[uid]
	return item, exists, nil
}
func (s *UserManagementStub) List(ctx context.Context) ([]User, error) {
	items := []User{}
	for _, i := range s.Users {
		items = append(items, i)
	}
	return items, nil
}

func (s *UserManagementStub) Query(ctx context.Context, filterName string) ([]User, error) {
	filterFunc, exists := s.filterService.GetUserFilterByName(ctx, filterName)
	if !exists {
		return nil, fmt.Errorf("User filter with name %s ot found", filterName)
	}
	users := []User{}
	for _, u := range s.Users {
		match, _ := filterFunc(ctx, u)
		if match {
			users = append(users, u)
		}
	}
	return users, nil
}
