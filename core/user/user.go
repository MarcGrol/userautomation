package user

import (
	"context"
	"log"
)

type User struct {
	UID        string
	Attributes map[string]interface{}
}

func (u User) HasAttributes(attributes []string) bool {
	for _, attr := range attributes {
		_, exists := u.Attributes[attr]
		if !exists {
			log.Printf("Missig mandatory attribute %s", attr)
			return false
		}
	}
	return true
}

type FilterFunc func(ctx context.Context, u User) (bool, error)

type Management interface {
	Put(ctx context.Context, user User) error
	Get(ctx context.Context, userUID string) (User, bool, error)
	List(ctx context.Context) ([]User, error)
	Query(ctx context.Context, filterName string) ([]User, error) // Could be a WHERE clause in the future
	Remove(ctx context.Context, userUID string) error
}

type UserFilterResolver interface {
	GetUserFilterByName(ctx context.Context, name string) (FilterFunc, bool)
}
