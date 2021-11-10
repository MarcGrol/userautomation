package userrule

import (
	"context"
	"fmt"
)

const (
	TriggerTopicName = "userruletrigger"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type EventHandler interface {
	OnRuleExecutionRequestedEvent(ctx context.Context, event RuleExecutionRequestedEvent) error
}

type RuleExecutionRequestedEvent struct {
	Rule Spec
}

func DispatchTriggerEvent(ctx context.Context, handler EventHandler, topic string, event interface{}) error {
	if topic != TriggerTopicName {
		return fmt.Errorf("Topic '%+v' is not right for user events. Must be: '%s'", topic, TriggerTopicName)
	}
	switch e := event.(type) {
	case RuleExecutionRequestedEvent:
		return handler.OnRuleExecutionRequestedEvent(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, TriggerTopicName)
	}
}
