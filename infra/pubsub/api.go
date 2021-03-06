package pubsub

import (
	"context"

	"github.com/gorilla/mux"
)

type OnEventFunc func(ctx context.Context, topic string, event interface{}) error

//go:generate mockgen -source=api.go -destination=pubsub_mock.go -package=pubsub Pubsub
type Pubsub interface {
	Subscribe(ctx context.Context, who, topic string, onEvent OnEventFunc) error
	Publish(ctx context.Context, topic string, event interface{}) error
}

// This interface is a marker that indicates that service is an event consumer
type SubscribingService interface {
	IamSubscribing()
	Subscribe(c context.Context, router *mux.Router) error
	OnEvent(c context.Context, topic string, event interface{}) error
}
