package segmentrule

import (
	"context"
	"fmt"
)

const (
	ManagementTopicName = "segmentrule"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type EventHandler interface {
	OnRuleCreated(ctx context.Context, event CreatedEvent) error
	OnRuleModified(ctx context.Context, event ModifiedEvent) error
	OnRuleRemoved(ctx context.Context, event RemovedEvent) error
}

type CreatedEvent struct {
	RuleState Spec
}

type ModifiedEvent struct {
	OldRuleState Spec
	NewRuleState Spec
}

type RemovedEvent struct {
	SegmentState Spec
}

func DispatchManagementEvent(ctx context.Context, handler EventHandler, topic string, event interface{}) error {
	if topic != ManagementTopicName {
		return fmt.Errorf("Topic '%+v' is not right for user events. Must be: '%s'", topic, ManagementTopicName)
	}
	switch e := event.(type) {
	case CreatedEvent:
		return handler.OnRuleCreated(ctx, e)
	case ModifiedEvent:
		return handler.OnRuleModified(ctx, e)
	case RemovedEvent:
		return handler.OnRuleRemoved(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, ManagementTopicName)
	}
}
