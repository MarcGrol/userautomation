package taskqueue

import (
	"context"
)

// TODO integrate with 3rd party product like Google pubsub

type simplisticTaskQueue struct {
}

func NewTaskQueue() TaskQueue {
	return &simplisticTaskQueue{}
}
func (tq *simplisticTaskQueue) Enqueue(ctx context.Context, task Task) error {

	return nil
}
