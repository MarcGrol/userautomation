package rule

import (
	"context"
	"fmt"
)

const (
	RuleTopicName = "rule"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type EventHandler interface {
	OnRuleCreated(ctx context.Context, rule UserSegmentRule) error
	OnRuleModified(ctx context.Context, oldState UserSegmentRule, newState UserSegmentRule) error
	OnRuleRemoved(ctx context.Context, rule UserSegmentRule) error
}

type CreatedEvent struct {
	State UserSegmentRule
}

type ModifiedEvent struct {
	OldState UserSegmentRule
	NewState UserSegmentRule
}

type RemovedEvent struct {
	State UserSegmentRule
}

func DispatchEvent(ctx context.Context, handler EventHandler, topic string, event interface{}) error {
	if topic != RuleTopicName {
		return fmt.Errorf("Topic '%+v' is not right for rule events. Must be: '%s'", topic, RuleTopicName)
	}
	switch e := event.(type) {
	case CreatedEvent:
		return handler.OnRuleCreated(ctx, e.State)
	case ModifiedEvent:
		return handler.OnRuleModified(ctx, e.OldState, e.NewState)
	case RemovedEvent:
		return handler.OnRuleRemoved(ctx, e.State)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, RuleTopicName)
	}
}
