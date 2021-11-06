package pubsub

import "context"

// This interface is a marker that indicates that service is an event consumer
type SubscribingService interface {
	Subscribe(c context.Context) error
	OnEvent(c context.Context, topic string, event interface{}) error
}
