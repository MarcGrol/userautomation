package rulemanagement

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/infra/datastore"
)

type ruleService struct {
	segmentStore datastore.Datastore
}

func New(segmentStore datastore.Datastore) rule.RuleService {
	return &ruleService{
		segmentStore: segmentStore,
	}
}

func (s *ruleService) Put(ctx context.Context, rule rule.RuleSpec) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		return s.segmentStore.Put(ctx, rule.UID, rule)
	})
}

func (s *ruleService) Get(ctx context.Context, ruleUID string) (rule.RuleSpec, bool, error) {
	r := rule.RuleSpec{}
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
		r = item.(rule.RuleSpec)
		ruleExists = exists

		return nil
	})
	if err != nil {
		return r, false, err
	}

	return r, ruleExists, nil
}

func (s *ruleService) Remove(ctx context.Context, ruleUID string) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		return s.segmentStore.Remove(ctx, ruleUID)
	})
}

func (s *ruleService) List(ctx context.Context) ([]rule.RuleSpec, error) {
	rules := []rule.RuleSpec{}
	var err error

	err = s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all rules: %s", err)
		}

		for _, item := range items {
			rules = append(rules, item.(rule.RuleSpec))
		}
		return nil
	})
	if err != nil {
		return rules, err
	}

	return rules, nil
}
