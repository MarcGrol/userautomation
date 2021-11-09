package usertaskexecutor

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/actions/emailaction"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/services/actionmanager"
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
	actionSpec := event.Task.RuleSpec.ActionSpec
	switch actionSpec.Name {
	case actionmanager.MailToOldName:
		return emailaction.NewEmailAction(
			actionSpec.ProvidedAttributes["email_subject"],
			actionSpec.ProvidedAttributes["email_body"],
			emailsending.NewEmailSender()).Perform(ctx, event.Task)
	case actionmanager.SmsToYoungName:
		return emailaction.NewEmailAction(
			actionSpec.ProvidedAttributes["email_subject"],
			actionSpec.ProvidedAttributes["email_body"],
			emailsending.NewEmailSender()).Perform(ctx, event.Task)
	default:
		return fmt.Errorf("Action %s not recoognized", actionSpec.Name)
	}
}
