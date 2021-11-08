package ondemandexecutor

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type onDemandRuleExecutor struct {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	rule.TriggerEventHandler
	pubsub      pubsub.Pubsub
	ruleService rule.SegmentRuleService
	userService user.Management
}

func New(pubsub pubsub.Pubsub, ruleService rule.SegmentRuleService, userService user.Management) rule.TriggerEventHandler {
	return &onDemandRuleExecutor{
		pubsub:      pubsub,
		ruleService: ruleService,
		userService: userService,
	}
}

type SegmentEventService interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	rule.TriggerEventHandler
}

func (s *onDemandRuleExecutor) IamSubscribing() {}

func (s *onDemandRuleExecutor) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, rule.TriggerTopicName, s.OnEvent)
}

func (s *onDemandRuleExecutor) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return rule.DispatchTriggerEvent(ctx, s, topic, event)
}

func (s *onDemandRuleExecutor) OnRuleExecutionRequestedEvent(ctx context.Context, event rule.RuleExecutionRequestedEvent) error {
	r, exists, err := s.ruleService.Get(ctx, event.Rule.UID)
	if err != nil {
		return fmt.Errorf("Error getting rule with uid %s: %s", event.Rule.UID, err)
	}
	if !exists {
		return fmt.Errorf("Rule with uid %s does not exist: %s", event.Rule.UID, err)
	}

	users, err := s.userService.QueryByName(ctx, r.UserSegment.UserFilterName)
	if err != nil {
		return fmt.Errorf("Error querying users: %s", err)
	}

	for _, u := range users {
		err = s.publishActionForUser(ctx, r, u)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *onDemandRuleExecutor) publishActionForUser(ctx context.Context, r rule.UserSegmentRule, user user.User) error {
	err := s.pubsub.Publish(ctx, action.TopicName, action.ActionExecutionRequestedEvent{
		Action: action.UserAction{
			RuleUID:  r.UID,
			Reason:   action.ReasonIsOnDemand,
			OldState: nil,
			NewState: &user,
		},
	})
	if err != nil {
		return fmt.Errorf("Error publishing action for rule %s and user %s: %s", r.UID, user.UID, err)
	}

	return nil
}
