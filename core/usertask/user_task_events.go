package usertask

import (
	"context"
	"fmt"
)

const (
	TopicName = "usertask"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type UserTaskEventHandler interface {
	OnUserTaskExecutionRequestedEvent(ctx context.Context, event UserTaskExecutionRequestedEvent) error
}

type UserTaskExecutionRequestedEvent struct {
	Task Spec
}

func DispatchEvent(ctx context.Context, handler UserTaskEventHandler, topic string, event interface{}) error {
	if topic != TopicName {
		return fmt.Errorf("Topic '%+v' is not right for user events. Must be: '%s'", topic, TopicName)
	}
	switch e := event.(type) {
	case UserTaskExecutionRequestedEvent:
		return handler.OnUserTaskExecutionRequestedEvent(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, TopicName)
	}
}
