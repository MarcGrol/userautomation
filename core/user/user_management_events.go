package user

import (
	"context"
	"fmt"
)

const (
	ManagementTopicName = "usermanagement"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type EventHandler interface {
	OnUserCreated(ctx context.Context, event CreatedEvent) error
	OnUserModified(ctx context.Context, event ModifiedEvent) error
	OnUserRemoved(ctx context.Context, event RemovedEvent) error
}

type CreatedEvent struct {
	UserState User
}

type ModifiedEvent struct {
	OldUserState User
	NewUserState User
}

type RemovedEvent struct {
	UserState User
}

func DispatchEvent(ctx context.Context, handler EventHandler, topic string, event interface{}) error {
	if topic != ManagementTopicName {
		return fmt.Errorf("Topic '%+v' is not right for user events. Must be: '%s'", topic, ManagementTopicName)
	}
	switch e := event.(type) {
	case CreatedEvent:
		return handler.OnUserCreated(ctx, e)
	case ModifiedEvent:
		return handler.OnUserModified(ctx, e)
	case RemovedEvent:
		return handler.OnUserRemoved(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, ManagementTopicName)
	}
}
