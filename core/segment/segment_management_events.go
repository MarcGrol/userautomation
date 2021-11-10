package segment

import (
	"context"
	"fmt"
)

const (
	ManagementTopicName = "segmentmanagement"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type EventHandler interface {
	OnSegmentCreated(ctx context.Context, event CreatedEvent) error
	OnSegmentModified(ctx context.Context, event ModifiedEvent) error
	OnSegmentRemoved(ctx context.Context, event RemovedEvent) error
}

type CreatedEvent struct {
	SegmentState Spec
}

type ModifiedEvent struct {
	OldSegmentState Spec
	NewSegmentState Spec
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
		return handler.OnSegmentCreated(ctx, e)
	case ModifiedEvent:
		return handler.OnSegmentModified(ctx, e)
	case RemovedEvent:
		return handler.OnSegmentRemoved(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, UserTopicName)
	}
}
