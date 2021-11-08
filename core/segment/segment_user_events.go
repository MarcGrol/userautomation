package segment

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/user"
)

const (
	UserTopicName = "usersegment"
)

// When new events are being introduced, this interface (and the dispatcher below) must be extended.
// Subscribers should implement this interface. This way, all subscribers can be easily detected and fixed (=extended)
type UserEventHandler interface {
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

func DispatchUserEvent(ctx context.Context, handler UserEventHandler, topic string, event interface{}) error {
	if topic != UserTopicName {
		return fmt.Errorf("Topic '%+v' is not right for segment events. Must be: '%s'", topic, UserTopicName)
	}
	switch e := event.(type) {
	case UserAddedToSegmentEvent:
		return handler.OnUserAddedToSegment(ctx, e)
	case UserRemovedFromSegmentEvent:
		return handler.OnUserRemovedFromSegment(ctx, e)
	default:
		return fmt.Errorf("Event %+v is not supported for topic %s", e, UserTopicName)
	}
}
