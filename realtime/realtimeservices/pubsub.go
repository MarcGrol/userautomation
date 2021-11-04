package realtimeservices

import (
	"context"
	"sync"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

type pubsub struct {
	sync.Mutex
	topics map[string][]realtimecore.OnEventFunc
}

func NewPubSub() realtimecore.Pubsub {
	return &pubsub{
		topics:map[string][]realtimecore.OnEventFunc{},
	}
}

func (ps *pubsub) Subscribe(ctx context.Context, topic string, onEvent realtimecore.OnEventFunc) error{
	ps.Lock()
	defer ps.Unlock()

	handlers, exists := ps.topics[topic]
	if !exists {
		handlers = []realtimecore.OnEventFunc{}
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



