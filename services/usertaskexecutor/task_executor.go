package usertaskexecutor

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type UserTaskExecutor interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	usertask.UserTaskEventHandler
}

type userTaskExecutor struct {
	pubsub pubsub.Pubsub
}

func New(pubsub pubsub.Pubsub) UserTaskExecutor {
	return &userTaskExecutor{
		pubsub: pubsub,
	}
}

func (s *userTaskExecutor) IamSubscribing() {}


func (s *userTaskExecutor) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, usertask.TopicName, s.OnEvent)
}

func (s *userTaskExecutor) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return usertask.DispatchEvent(ctx, s, topic, event)
}

func (s *userTaskExecutor) OnUserTaskExecutionRequestedEvent(ctx context.Context, event usertask.UserTaskExecutionRequestedEvent) error {
	return fmt.Errorf("Not implemented")
}
