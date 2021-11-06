package pubsub

import (
	"context"
	"sync"
)

type pubsub struct {
	sync.Mutex
	topics map[string][]OnEventFunc
}

func NewPubSub() Pubsub {
	return &pubsub{
		topics: map[string][]OnEventFunc{},
	}
}

func (ps *pubsub) Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error {
	ps.Lock()
	defer ps.Unlock()

	handlers, exists := ps.topics[topic]
	if !exists {
		handlers = []OnEventFunc{}
	}
	handlers = append(handlers, onEvent)
	ps.topics[topic] = handlers

	return nil
}

func (ps *pubsub) Publish(ctx context.Context, topic string, event interface{}) error {
	ps.Lock()
	defer ps.Unlock()

	handlers, found := ps.topics[topic]
	if !found {
		return nil
	}

	for _, handler := range handlers {
		err := handler(ctx, topic, event)
		if err != nil {
			return err
		}
	}

	return nil
}
