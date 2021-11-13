package pubsub

import (
	"context"
	"log"
)

// TODO integrate with 3rd party product like Google pubsub

type simplisticPubsub struct {
	topics map[string][]subscription
}

func NewPubSub() Pubsub {
	return &simplisticPubsub{
		topics: map[string][]subscription{},
	}
}

type subscription struct {
	who     string
	topic   string
	onEvent OnEventFunc
}

func (ps *simplisticPubsub) Subscribe(ctx context.Context, who, topic string, onEvent OnEventFunc) error {
	subscriptions, exists := ps.topics[topic]
	if !exists {
		subscriptions = []subscription{}
	}
	subscriptions = append(subscriptions, subscription{
		who:     who,
		topic:   topic,
		onEvent: onEvent,
	})
	ps.topics[topic] = subscriptions

	return nil
}

func (ps *simplisticPubsub) Publish(ctx context.Context, topic string, event interface{}) error {
	subscriptions, found := ps.topics[topic]
	if !found {
		return nil
	}

	for _, subscription := range subscriptions {
		log.Printf("Publish to %s using topic %s: %+v", subscription.who, topic, event)
		// This should run in the background
		err := subscription.onEvent(ctx, topic, event)
		if err != nil {
			log.Printf("Error handling event: %s: %+v", topic, event)
			// Although this single subscriber fails handling the event,
			// all the other subscribers should still be triggered
		}
	}

	return nil
}
