package rules

import (
	"context"
	"github.com/MarcGrol/userautomation/users"
)

type UserChangeStatus int

const (
	UserCreated UserChangeStatus = iota
	UserModified
	UserRemoved
)

type UserActionFunc func(ctx context.Context, ruleName string, changeStatus UserChangeStatus, oldState *users.User, newState *users.User) error

type UserSegmentRule struct {
	Name                string
	IsApplicableForUser users.UserFilterFunc // Could use a WHERE clause alternatively
	PerformAction       UserActionFunc
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, name string) (UserSegmentRule, bool, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}
