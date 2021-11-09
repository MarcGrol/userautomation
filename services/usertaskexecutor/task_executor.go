package usertaskexecutor

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/infra/taskqueue"
)

type UserTaskExecutor interface {
	// Flags that this service is an event consumer
	taskqueue.TaskQueueReceiver
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
}

type userTaskExecutor struct {
	taskqueue taskqueue.TaskQueue
}

func New(taskqueue taskqueue.TaskQueue) UserTaskExecutor {
	return &userTaskExecutor{
		taskqueue: taskqueue,
	}
}
func (s *userTaskExecutor) IamReceivingTasks() {}

func (s *userTaskExecutor) OnTaskReceived(c context.Context, queueName string, payload string) error {
	return fmt.Errorf("Not implemented")
}
