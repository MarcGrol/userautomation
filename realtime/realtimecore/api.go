package realtimecore

import "context"

type OnEventFunc func( ctx context.Context, topic string, event interface{}) error
type Pubsub interface {
	Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error
	Publish(ctx context.Context, topic string, event interface{}) error
}

type SegmentRuleService interface {
	Put(ctx context.Context, segmentRule UserSegmentRule) error
	Get(ctx context.Context, name string) (UserSegmentRule, bool, error)
	Delete(ctx context.Context, name string) error
	List(ctx context.Context) ([]UserSegmentRule, error)
}

type UserService interface {
	Put(ctx context.Context, user User) error
	Get(ctx context.Context, userUID string) (User, bool, error)
	Query(ctx context.Context, filter UserFilterFunc) ([]User, error) // Could use a WHERE clause alternatively
	Delete(ctx context.Context, userUID string) error
}

type UserEventService interface {
	Subscribe(c context.Context) error
	OnEvent(c context.Context, topic string, event interface{}) error
}
