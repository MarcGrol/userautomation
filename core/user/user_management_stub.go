package user

import (
	"context"
	"fmt"
)

type UserManagementStub struct {
	Users map[string]User
}

func NewUserManagementStub() *UserManagementStub {
	return &UserManagementStub{
		Users: map[string]User{},
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

func (s *UserManagementStub) QueryByName(ctx context.Context, filterName string) ([]User, error) {
	filterFunc, found := GetUserFilterByName(ctx, filterName)
	if !found {
		return []User{}, fmt.Errorf("Filter func with name %s does not exist", filterName)
	}
	return s.QueryByFunc(ctx, filterFunc)
}

func (s *UserManagementStub) QueryByFunc(ctx context.Context, filter FilterFunc) ([]User, error) {
	users := []User{}
	for _, u := range s.Users {
		match, _ := filter(ctx, u)
		if match {
			users = append(users, u)
		}
	}
	return users, nil
}
