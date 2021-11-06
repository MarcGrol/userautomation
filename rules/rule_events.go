package rules

import (
	"context"
	"fmt"
)

const (
	RuleTopicName = "rule"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type RuleEventHandler interface {
	OnRuleCreated(ctx context.Context, rule UserSegmentRule) error
	OnRuleModified(ctx context.Context, oldState UserSegmentRule, newState UserSegmentRule) error
	OnRuleRemoved(ctx context.Context, rule UserSegmentRule) error
}

type RuleCreatedEvent struct {
	State UserSegmentRule
}

type RuleModifiedEvent struct {
	OldState UserSegmentRule
	NewState UserSegmentRule
}

type RuleRemovedEvent struct {
	State UserSegmentRule
}

func DispatchEvent(ctx context.Context, handler RuleEventHandler, topic string, event interface{}) error {
	if topic != RuleTopicName {
		return fmt.Errorf("Topic '%+v' is not right for rule events. Must be: '%s'", topic, RuleTopicName)
	}
	switch e := event.(type) {
	case RuleCreatedEvent:
		return handler.OnRuleCreated(ctx, e.State)
	case RuleModifiedEvent:
		return handler.OnRuleModified(ctx, e.OldState, e.NewState)
	case RuleRemovedEvent:
		return handler.OnRuleRemoved(ctx, e.State)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, RuleTopicName)
	}
}
