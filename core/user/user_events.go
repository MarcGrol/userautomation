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
type EventHandler interface {
	OnUserCreated(ctx context.Context, user User) error
	OnUserModified(ctx context.Context, oldState User, newState User) error
	OnUserRemoved(ctx context.Context, user User) error
}

type CreatedEvent struct {
	State User
}

type ModifiedEvent struct {
	OldState User
	NewState User
}

type RemovedEvent struct {
	State User
}

func DispatchEvent(ctx context.Context, handler EventHandler, topic string, event interface{}) error {
	if topic != UserTopicName {
		return fmt.Errorf("Topic '%+v' is not right for user events. Must be: '%s'", topic, UserTopicName)
	}
	switch e := event.(type) {
	case CreatedEvent:
		return handler.OnUserCreated(ctx, e.State)
	case ModifiedEvent:
		return handler.OnUserModified(ctx, e.OldState, e.NewState)
	case RemovedEvent:
		return handler.OnUserRemoved(ctx, e.State)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, UserTopicName)
	}
}
