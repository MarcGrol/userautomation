package user

import (
	"context"
	"fmt"
)

type UserServiceStub struct {
	Users map[string]User
}

func NewUserServiceStub() *UserServiceStub {
	return &UserServiceStub{
		Users: map[string]User{},
	}
}

func (s *UserServiceStub) Put(ctx context.Context, u User) error {
	s.Users[u.UID] = u
	return nil
}

func (s *UserServiceStub) Remove(ctx context.Context, uid string) error {
	delete(s.Users, uid)
	return nil
}

func (s *UserServiceStub) Get(ctx context.Context, uid string) (User, bool, error) {
	item, exists := s.Users[uid]
	return item, exists, nil
}
func (s *UserServiceStub) List(ctx context.Context) ([]User, error) {
	items := []User{}
	for _, i := range s.Users {
		items = append(items, i)
	}
	return items, nil
}

func (s *UserServiceStub) QueryByName(ctx context.Context, filterName string) ([]User, error) {
	filterFunc, found := GetUserFilterByName(ctx, filterName)
	if !found {
		return []User{}, fmt.Errorf("Filter func with name %s does not exist", filterName)
	}
	return s.QueryByFunc(ctx, filterFunc)
}

func (s *UserServiceStub) QueryByFunc(ctx context.Context, filter FilterFunc) ([]User, error) {
	users := []User{}
	for _, u := range s.Users {
		match, _ := filter(ctx, u)
		if match {
			users = append(users, u)
		}
	}
	return users, nil
}
