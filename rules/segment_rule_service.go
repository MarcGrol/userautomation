package rules

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/infra/pubsub"

	"github.com/MarcGrol/userautomation/infra/datastore"
)

type userSegmentRuleService struct {
	datastore datastore.Datastore
	pubsub    pubsub.Pubsub
}

func NewUserSegmentRuleService(datastore datastore.Datastore, pubsub pubsub.Pubsub) SegmentRuleService {
	return &userSegmentRuleService{
		datastore: datastore,
		pubsub:    pubsub, // TODO might signal rule change to dedicated service
	}
}

func (s *userSegmentRuleService) Put(ctx context.Context, segmentRule UserSegmentRule) error {
	return s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		return s.datastore.Put(ctx, segmentRule.Name, segmentRule)
	})
}

func (s *userSegmentRuleService) Get(ctx context.Context, ruleName string) (UserSegmentRule, bool, error) {
	rule := UserSegmentRule{}
	ruleExists := false
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.datastore.Get(ctx, ruleName)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", ruleName, err)
		}
		rule = item.(UserSegmentRule)
		ruleExists = exists

		return nil
	})
	if err != nil {
		return rule, false, err
	}

	return rule, ruleExists, nil
}

func (s *userSegmentRuleService) Delete(ctx context.Context, ruleName string) error {
	return s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		return s.datastore.Remove(ctx, ruleName)
	})
}

func (s *userSegmentRuleService) List(ctx context.Context) ([]UserSegmentRule, error) {
	rules := []UserSegmentRule{}
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.datastore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all rules: %s", err)
		}

		for _, item := range items {
			rules = append(rules, item.(UserSegmentRule))
		}
		return nil
	})
	if err != nil {
		return rules, err
	}

	return rules, nil
}
