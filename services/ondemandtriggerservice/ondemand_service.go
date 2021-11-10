package ondemandtriggerservice

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type OnDemandService interface {
	segmentrule.TriggerRuleExecution
}

type onDemandService struct {
	ruleService segmentrule.Service
	pubsub      pubsub.Pubsub
}

func New(ruleService segmentrule.Service, pubsub pubsub.Pubsub) OnDemandService {
	return &onDemandService{
		ruleService: ruleService,
		pubsub:      pubsub,
	}
}

func (s *onDemandService) Trigger(ctx context.Context, ruleUID string) error {
	r, exists, err := s.ruleService.Get(ctx, ruleUID)
	if err != nil {
		return fmt.Errorf("Error getting rule with uid %s: %s", ruleUID, err)
	}
	if !exists {
		return fmt.Errorf("Rule with uid %s does not exist: %s", ruleUID, err)
	}

	err = s.pubsub.Publish(ctx, segmentrule.TriggerTopicName, segmentrule.RuleExecutionRequestedEvent{Rule: r})
	if err != nil {
		return fmt.Errorf("Error publishing RuleExecutionRequestedEvent for rule %+v: %s", r, err)
	}

	return nil
}
