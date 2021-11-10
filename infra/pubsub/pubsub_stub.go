package pubsub

import "context"

type Publication struct {
	Topic string
	Event interface{}
}

type PubsubStub struct {
	Publications []Publication
}

func NewPubsubStub() *PubsubStub {
	return &PubsubStub{
		Publications: []Publication{},
	}
}

func (ps *PubsubStub) Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error {
	return nil
}

func (ps *PubsubStub) Publish(ctx context.Context, topic string, event interface{}) error {
	ps.Publications = append(ps.Publications, Publication{
		Topic: topic,
		Event: event,
	})
	return nil
}
