package segment

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/user"
)

const (
	TopicName = "segment"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type EventHandler interface {
	OnUserAddedToSegment(ctx context.Context, event UserAddedToSegmentEvent) error
	OnUserRemovedFromSegment(ctx context.Context, event UserRemovedFromSegmentEvent) error
}

type UserAddedToSegmentEvent struct {
	SegmentUID string
	User       user.User
}

type UserRemovedFromSegmentEvent struct {
	SegmentUID string
	User       user.User
}

func DispatchEvent(ctx context.Context, handler EventHandler, topic string, event interface{}) error {
	if topic != TopicName {
		return fmt.Errorf("Topic '%+v' is not right for segment events. Must be: '%s'", topic, TopicName)
	}
	switch e := event.(type) {
	case UserAddedToSegmentEvent:
		return handler.OnUserAddedToSegment(ctx, e)
	case UserRemovedFromSegmentEvent:
		return handler.OnUserRemovedFromSegment(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, TopicName)
	}
}
