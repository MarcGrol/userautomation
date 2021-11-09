package ondemandtriggerservice

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type OnDemandService struct {
	ruleService rule.RuleService
	pubsub      pubsub.Pubsub
}

func New(ruleService rule.RuleService, pubsub pubsub.Pubsub) rule.TriggerRuleExecution {
	return &OnDemandService{
		ruleService: ruleService,
		pubsub:      pubsub,
	}
}

func (s *OnDemandService) Trigger(ctx context.Context, ruleUID string) error {
	r, exists, err := s.ruleService.Get(ctx, ruleUID)
	if err != nil {
		return fmt.Errorf("Error getting rule with uid %s: %s", ruleUID, err)
	}
	if !exists {
		return fmt.Errorf("Rule with uid %s does not exist: %s", ruleUID, err)
	}

	err = s.pubsub.Publish(ctx, rule.TriggerTopicName, rule.RuleExecutionRequestedEvent{Rule: r})
	if err != nil {
		return fmt.Errorf("Error publishing RuleExecutionRequestedEvent for rule %+v: %s", r, err)
	}

	return nil
}
