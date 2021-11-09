package taskqueue

import "context"

type TaskQueueStub struct {
	TasksEnqueued []Task
}

func NewTaskQueueStub() TaskQueue {
	return &TaskQueueStub{
		TasksEnqueued: []Task{},
	}
}
func (tq *TaskQueueStub) Enqueue(ctx context.Context, task Task) error {
	tq.TasksEnqueued = append(tq.TasksEnqueued, task)
	return nil
}
