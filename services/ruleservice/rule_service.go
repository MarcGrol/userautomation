package ruleservice

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/infra/datastore"
)

type userSegmentRuleService struct {
	segmentStore datastore.Datastore
}

func NewUserSegmentRuleService(segmentStore datastore.Datastore) rule.SegmentRuleService {
	return &userSegmentRuleService{
		segmentStore: segmentStore,
	}
}

func (s *userSegmentRuleService) Put(ctx context.Context, rule rule.UserSegmentRule) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		return s.segmentStore.Put(ctx, rule.UID, rule)
	})
}

func (s *userSegmentRuleService) Get(ctx context.Context, ruleUID string) (rule.UserSegmentRule, bool, error) {
	r := rule.UserSegmentRule{}
	ruleExists := false
	var err error

	err = s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.segmentStore.Get(ctx, ruleUID)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", ruleUID, err)
		}
		if !exists {
			return fmt.Errorf("Rule with uid %s not found", ruleUID)
		}
		r = item.(rule.UserSegmentRule)
		ruleExists = exists

		return nil
	})
	if err != nil {
		return r, false, err
	}

	return r, ruleExists, nil
}

func (s *userSegmentRuleService) Remove(ctx context.Context, ruleUID string) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		return s.segmentStore.Remove(ctx, ruleUID)
	})
}

func (s *userSegmentRuleService) List(ctx context.Context) ([]rule.UserSegmentRule, error) {
	rules := []rule.UserSegmentRule{}
	var err error

	err = s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all rules: %s", err)
		}

		for _, item := range items {
			rules = append(rules, item.(rule.UserSegmentRule))
		}
		return nil
	})
	if err != nil {
		return rules, err
	}

	return rules, nil
}
