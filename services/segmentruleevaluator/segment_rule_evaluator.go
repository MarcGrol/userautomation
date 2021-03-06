package segmentruleevaluator

import (
	"context"
	"fmt"
	"log"

	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/gorilla/mux"
)

type Service interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	segmentrule.TriggerEventHandler
}

type service struct {
	pubsub      pubsub.Pubsub
	ruleService segmentrule.Management
	userService user.Management
}

func New(pubsub pubsub.Pubsub, ruleService segmentrule.Management, userService user.Management) Service {
	return &service{
		pubsub:      pubsub,
		ruleService: ruleService,
		userService: userService,
	}
}

func (s *service) IamSubscribing() {}

func (s *service) Subscribe(ctx context.Context, router *mux.Router) error {
	return s.pubsub.Subscribe(ctx, "segmentruleevaluator", segmentrule.TriggerTopicName, s.OnEvent)
}

func (s *service) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return segmentrule.DispatchTriggerEvent(ctx, s, topic, event)
}

func (s *service) OnRuleExecutionRequestedEvent(ctx context.Context, event segmentrule.RuleExecutionRequestedEvent) error {
	rule, exists, err := s.ruleService.Get(ctx, event.Rule.UID)
	if err != nil {
		return fmt.Errorf("Error getting rule with uid %s: %s", event.Rule.UID, err)
	}
	if !exists {
		return fmt.Errorf("Rule with uid %s does not exist: %s", event.Rule.UID, err)
	}

	log.Printf("Start evaluating rule %+v", rule)

	// TODO this possibly a very large task that would lock the datastore for a long time:
	// we might want to break this up with cursors into multiple tasks
	users, err := s.userService.Query(ctx, rule.SegmentSpec.UserFilterName)
	if err != nil {
		return fmt.Errorf("Error querying users: %s", err)
	}
	log.Printf("Found %d users: %+v", len(users), users)

	for _, u := range users {
		log.Printf("Evaluate user: %+v", u)
		if !u.HasAttributes(event.Rule.ActionSpec.MandatoryUserAttributes) {
			return fmt.Errorf("User %s is missing madatory attributes for action %s", u.UID, event.Rule.ActionSpec.Name)
		}
		err = s.publishActionForUser(ctx, rule, u)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *service) publishActionForUser(ctx context.Context, r segmentrule.Spec, u user.User) error {
	err := s.pubsub.Publish(ctx, usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
		Task: usertask.Spec{
			UID:        "", // TODO identify each triggered rule uninquely
			RuleUID:    r.UID,
			ActionSpec: r.ActionSpec,
			Reason:     usertask.ReasonSegmentRuleTriggered,
			User:       u,
		},
	})
	if err != nil {
		return fmt.Errorf("Error publishing task for rule %s and user %s: %s", r.UID, u.UID, err)
	}

	return nil
}
