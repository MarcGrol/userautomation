package pubsub

import "context"

type PubsubStub struct{}

func NewPubsubStub() *PubsubStub {
	return &PubsubStub{}
}

func (ps PubsubStub) Subscribe(ctx context.Context, topic string, onEvent OnEventFunc) error {
	return nil
}

func (ps PubsubStub) Publish(ctx context.Context, topic string, event interface{}) error {
	return nil
}
