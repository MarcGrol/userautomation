package userruletriggerservice

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/userrule"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/gorilla/mux"
)

type service struct {
	pubsub pubsub.Pubsub
}

func New(pubsub pubsub.Pubsub) userrule.TriggerRuleExecution {
	return &service{
		pubsub: pubsub,
	}
}

func (s *service) Trigger(ctx context.Context, rule userrule.Spec) error {
	// TODO validate rule
	err := s.pubsub.Publish(ctx, userrule.TriggerTopicName, userrule.RuleExecutionRequestedEvent{Rule: rule})
	if err != nil {
		return fmt.Errorf("Error publishing RuleExecutionRequestedEvent for rule %+v: %s", rule, err)
	}

	return nil
}

func (m *service) RegisterEndpoints(ctx context.Context, router *mux.Router) {

}
