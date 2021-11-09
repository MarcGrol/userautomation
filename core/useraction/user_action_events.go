package useraction

import (
	"context"
	"fmt"
)

const (
	TopicName = "action"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type ActionExecutor interface {
	OnActionExecutionRequestedEvent(ctx context.Context, event ActionExecutionRequestedEvent) error
}

type ActionExecutionRequestedEvent struct {
	Action UserAction
}

func DispatchEvent(ctx context.Context, handler ActionExecutor, topic string, event interface{}) error {
	if topic != TopicName {
		return fmt.Errorf("Topic '%+v' is not right for user events. Must be: '%s'", topic, TopicName)
	}
	switch e := event.(type) {
	case ActionExecutionRequestedEvent:
		return handler.OnActionExecutionRequestedEvent(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, TopicName)
	}
}
