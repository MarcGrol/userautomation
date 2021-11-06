package user

import (
	"context"
	"fmt"
)

const (
	UserTopicName = "user"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type UserEventHandler interface {
	OnUserCreated(ctx context.Context, user User) error
	OnUserModified(ctx context.Context, oldState User, newState User) error
	OnUserRemoved(ctx context.Context, user User) error
}

type UserCreatedEvent struct {
	State User
}

type UserModifiedEvent struct {
	OldState User
	NewState User
}

type UserRemovedEvent struct {
	State User
}

func DispatchEvent(ctx context.Context, handler UserEventHandler, topic string, event interface{}) error {
	if topic != UserTopicName {
		return fmt.Errorf("Topic '%+v' is not right for user events. Must be: '%s'", topic, UserTopicName)
	}
	switch e := event.(type) {
	case UserCreatedEvent:
		return handler.OnUserCreated(ctx, e.State)
	case UserModifiedEvent:
		return handler.OnUserModified(ctx, e.OldState, e.NewState)
	case UserRemovedEvent:
		return handler.OnUserRemoved(ctx, e.State)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, UserTopicName)
	}
}
