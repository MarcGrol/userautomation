package pubsub

import "context"

type OnEventFunc func(ctx context.Context, topic string, event interface{}) error

//go:generate mockgen -source=api.go -destination=pubsub_mock.go -package=pubsub Pubsub
type Pubsub interface {
	Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error
	Publish(ctx context.Context, topic string, event interface{}) error
}
