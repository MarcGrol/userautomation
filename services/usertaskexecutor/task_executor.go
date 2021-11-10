package usertaskexecutor

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/actions/emailaction"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
)

type Service interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	usertask.UserTaskEventHandler
}

type service struct {
	pubsub pubsub.Pubsub
}

func New(pubsub pubsub.Pubsub) Service {
	return &service{
		pubsub: pubsub,
	}
}

func (s *service) IamSubscribing() {}

func (s *service) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, usertask.TopicName, s.OnEvent)
}

func (s *service) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return usertask.DispatchEvent(ctx, s, topic, event)
}

func (s *service) OnUserTaskExecutionRequestedEvent(ctx context.Context, event usertask.UserTaskExecutionRequestedEvent) error {
	actionSpec := event.Task.ActionSpec
	switch actionSpec.Name {
	case supportedactions.MailToOldName:
		return emailaction.NewEmailAction(
			actionSpec.ProvidedInformation["email_subject"],
			actionSpec.ProvidedInformation["email_body"],
			emailsending.NewEmailSender()).Perform(ctx, event.Task)
	case supportedactions.SmsToYoungName:
		return emailaction.NewEmailAction(
			actionSpec.ProvidedInformation["email_subject"],
			actionSpec.ProvidedInformation["email_body"],
			emailsending.NewEmailSender()).Perform(ctx, event.Task)
	default:
		return fmt.Errorf("Action %s not recoognized", actionSpec.Name)
	}
}
