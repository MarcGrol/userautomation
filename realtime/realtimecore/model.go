package realtimecore

import "context"

type User struct {
	UID        string
	FullName   string
	Attributes map[string]interface{}
	Channels   []Channel
}

type Channel struct {
	Type    ChannelType
	Address string
}

type ChannelType int

const (
	ChannelTypeUnknown ChannelType = iota
	ChannelTypeEmail
	ChannelTypeSms
	ChannelTypeCustom
)

type UserStatus int

const (
	UserStatusUnknown UserStatus = iota
	UserCreated
	UserModified
	UserRemoved
)

type UserFilterFunc func(ctx context.Context, u User) (bool, error)
type UserActionFunc func(ctx context.Context, ruleName string, userStatus UserStatus, u User) error

type UserSegmentRule struct {
	Name                string
	IsApplicableForUser UserFilterFunc
	PerformAction       UserActionFunc
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, name string) (UserSegmentRule, bool, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}

type UserEventService interface {
	OnUserCreated(c context.Context, u User) error
	OnUserModified(c context.Context, before User, after User) error
	OnUserRemoved(c context.Context, u User) error
}

type UserService interface {
	Put(ctx context.Context, user User) error
	Get(ctx context.Context, userUID string) (User, bool, error)
	Query(ctx context.Context, filter UserFilterFunc) ([]User, error)
	Delete(ctx context.Context, userUID string) error
}
