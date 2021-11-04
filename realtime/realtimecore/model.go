package realtimecore

import "context"

type User struct {
	UID        string
	Attributes map[string]interface{}
}

type UserStatus int

const (
	UserCreated UserStatus = iota
	UserModified
	UserRemoved
)

type UserFilterFunc func(ctx context.Context, u User) (bool, error)
type UserActionFunc func(ctx context.Context, ruleName string, userStatus UserStatus, u User) error

type UserSegmentRule struct {
	Name                string
	IsApplicableForUser UserFilterFunc // Could use a WHERE clause alternatively
	PerformAction       UserActionFunc
}
