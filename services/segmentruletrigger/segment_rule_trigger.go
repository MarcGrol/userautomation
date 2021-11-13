package segmentruletrigger

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type Service interface {
	segmentrule.TriggerRuleExecution
}

type service struct {
	ruleService segmentrule.Management
	pubsub      pubsub.Pubsub
}

func New(ruleService segmentrule.Management, pubsub pubsub.Pubsub) Service {
	return &service{
		ruleService: ruleService,
		pubsub:      pubsub,
	}
}

func (s *service) Trigger(ctx context.Context, rule segmentrule.Spec) error {
	// TODO validate rule
	err := s.pubsub.Publish(ctx, segmentrule.TriggerTopicName, segmentrule.RuleExecutionRequestedEvent{Rule: rule})
	if err != nil {
		return fmt.Errorf("Error publishing RuleExecutionRequestedEvent for rule %+v: %s", rule, err)
	}

	return nil
}
