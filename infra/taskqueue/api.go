package taskqueue

import "context"

type Task struct {
	QueueName string
	Payload   string
}

//go:generate mockgen -source=api.go -destination=taskqueue_mock.go -package=taskqueue TaskQueue
type TaskQueue interface {
	Enqueue(ctx context.Context, task Task) error
}

// This interface is a marker that indicates that service is an task consumer
type TaskQueueReceiver interface {
	IamReceivingTasks()
	OnTaskReceived(c context.Context, queueName string, payload string) error
}
