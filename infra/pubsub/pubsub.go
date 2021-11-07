package pubsub

import (
	"context"
	"sync"
)

// TODO integrate with 3rd party product like Google pubsub

type simplisticPubsub struct {
	sync.Mutex
	topics map[string][]OnEventFunc
}

func NewPubSub() Pubsub {
	return &simplisticPubsub{
		topics: map[string][]OnEventFunc{},
	}
}

func (ps *simplisticPubsub) Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error {
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

func (ps *simplisticPubsub) Publish(ctx context.Context, topic string, event interface{}) error {
	ps.Lock()
	defer ps.Unlock()

	handlers, found := ps.topics[topic]
	if !found {
		return nil
	}

	for _, handler := range handlers {
		// This should run in the background
		err := handler(ctx, topic, event)
		if err != nil {
			// Although this single subscriber fails handling the event,
			// all the other subscribers should still be triggered
		}
	}

	return nil
}
