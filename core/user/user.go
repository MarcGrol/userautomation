package user

import (
	"context"
)

type User struct {
	UID        string
	Attributes map[string]interface{}
}

type FilterFunc func(ctx context.Context, u User) (bool, error)

type Service interface {
	Put(ctx context.Context, user User) error
	Get(ctx context.Context, userUID string) (User, bool, error)
	List(ctx context.Context) ([]User, error)
	QueryByFunc(ctx context.Context, filter FilterFunc) ([]User, error) // Could use a WHERE clause alternatively
	QueryByName(ctx context.Context, filterName string) ([]User, error) // Could use a WHERE clause alternatively
	Remove(ctx context.Context, userUID string) error
}
