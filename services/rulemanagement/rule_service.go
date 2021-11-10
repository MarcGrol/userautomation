package rulemanagement

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/infra/datastore"
)

type ruleService struct {
	segmentStore datastore.Datastore
}

func New(segmentStore datastore.Datastore) segmentrule.Service {
	return &ruleService{
		segmentStore: segmentStore,
	}
}

func (s *ruleService) Put(ctx context.Context, rule segmentrule.Spec) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		return s.segmentStore.Put(ctx, rule.UID, rule)
	})
}

func (s *ruleService) Get(ctx context.Context, ruleUID string) (segmentrule.Spec, bool, error) {
	r := segmentrule.Spec{}
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
		r = item.(segmentrule.Spec)
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

func (s *ruleService) List(ctx context.Context) ([]segmentrule.Spec, error) {
	rules := []segmentrule.Spec{}
	var err error

	err = s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all rules: %s", err)
		}

		for _, item := range items {
			rules = append(rules, item.(segmentrule.Spec))
		}
		return nil
	})
	if err != nil {
		return rules, err
	}

	return rules, nil
}
