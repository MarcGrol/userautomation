package pubsub

import "context"

type OnEventFunc func(ctx context.Context, topic string, event interface{}) error
type Pubsub interface {
	Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error
	Publish(ctx context.Context, topic string, event interface{}) error
}

