package usertaskexecutor

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/actions/emailaction"
	"github.com/MarcGrol/userautomation/actions/smsaction"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/integrations/smssending"
)

type Service interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	usertask.UserTaskEventHandler
}

type service struct {
	pubsub   pubsub.Pubsub
	reporter usertask.ExecutionReporter
}

func New(pubsub pubsub.Pubsub, reporter usertask.ExecutionReporter) Service {
	return &service{
		pubsub:   pubsub,
		reporter: reporter,
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
		{
			report, err := emailaction.NewEmailAction(
				actionSpec.ProvidedInformation["subject_template"],
				actionSpec.ProvidedInformation["body_template"],
				emailsending.NewEmailSender()).Perform(ctx, event.Task)
			if err != nil {
				return err
			}
			s.reporter.ReportExecution(ctx, report)
			return nil
		}
	case supportedactions.SmsToYoungName:
		{
			report, err := smsaction.New(
				actionSpec.ProvidedInformation["body_template"],
				smssending.NewSmsSender()).Perform(ctx, event.Task)
			if err != nil {
				return err
			}
			s.reporter.ReportExecution(ctx, report)
			return nil
		}
	default:
		return fmt.Errorf("Action %s not recoognized", actionSpec.Name)
	}
}
