package userruleevaluator

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/userrule"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type Service interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	userrule.TriggerEventHandler
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
	return s.pubsub.Subscribe(ctx, userrule.TriggerTopicName, s.OnEvent)
}

func (s *service) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return userrule.DispatchTriggerEvent(ctx, s, topic, event)
}

func (s *service) OnRuleExecutionRequestedEvent(ctx context.Context, event userrule.RuleExecutionRequestedEvent) error {
	u := event.Rule.User

	if !u.HasAttributes(event.Rule.ActionSpec.MandatoryUserAttributes) {
		return fmt.Errorf("User %s is missing madatory attributes for action %s", u.UID, event.Rule.ActionSpec.Name)
	}
	err := s.publishActionForUser(ctx, event.Rule)
	if err != nil {
		return err
	}

	return nil
}

func (s *service) publishActionForUser(ctx context.Context, rule userrule.Spec) error {
	err := s.pubsub.Publish(ctx, usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
		Task: usertask.Spec{
			RuleUID:    "", // No rule that triggered this task
			ActionSpec: rule.ActionSpec,
			Reason:     usertask.ReasonUserRuleExecuted,
			User:       rule.User,
		},
	})
	if err != nil {
		return fmt.Errorf("Error publishing task for rule %s and user %s: %s", rule.UID, rule.User.UID, err)
	}

	return nil
}
