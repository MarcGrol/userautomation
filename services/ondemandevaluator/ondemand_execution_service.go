package ondemandevaluator

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type OnDemandRuleEvaluator interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	rule.TriggerEventHandler
}

type onDemandRuleEvaluator struct {
	pubsub      pubsub.Pubsub
	ruleService rule.RuleService
	userService user.Management
}

func New(pubsub pubsub.Pubsub,  ruleService rule.RuleService, userService user.Management) OnDemandRuleEvaluator {
	return &onDemandRuleEvaluator{
		pubsub:      pubsub,
		ruleService: ruleService,
		userService: userService,
	}
}

func (s *onDemandRuleEvaluator) IamSubscribing() {}

func (s *onDemandRuleEvaluator) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, rule.TriggerTopicName, s.OnEvent)
}

func (s *onDemandRuleEvaluator) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return rule.DispatchTriggerEvent(ctx, s, topic, event)
}

func (s *onDemandRuleEvaluator) OnRuleExecutionRequestedEvent(ctx context.Context, event rule.RuleExecutionRequestedEvent) error {
	r, exists, err := s.ruleService.Get(ctx, event.Rule.UID)
	if err != nil {
		return fmt.Errorf("Error getting rule with uid %s: %s", event.Rule.UID, err)
	}
	if !exists {
		return fmt.Errorf("Rule with uid %s does not exist: %s", event.Rule.UID, err)
	}

	users, err := s.userService.QueryByName(ctx, r.SegmentSpec.UserFilterName)
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

func (s *onDemandRuleEvaluator) publishActionForUser(ctx context.Context, r rule.RuleSpec, u user.User) error {
	err := s.pubsub.Publish(ctx, usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
		Task: usertask.UserTask{
			RuleSpec: r,
			Reason:   usertask.ReasonIsOnDemand,
			User:     u,
		},
	})
	if err != nil {
		return fmt.Errorf("Error publishing task for rule %s and user %s: %s", r.UID, u.UID, err)
	}

	return nil
}
