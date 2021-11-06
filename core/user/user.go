package user

import (
	"context"
)

type User struct {
	UID        string
	Attributes map[string]interface{}
}

type UserFilterFunc func(ctx context.Context, u User) (bool, error)

type UserService interface {
	Put(ctx context.Context, user User) error
	Get(ctx context.Context, userUID string) (User, bool, error)
	Query(ctx context.Context, filter UserFilterFunc) ([]User, error) // Could use a WHERE clause alternatively
	Remove(ctx context.Context, userUID string) error
}
