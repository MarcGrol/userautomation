package rulemanagement

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/coredata/predefinedrules"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type service struct {
	ruleStore datastore.Datastore
	pubsub    pubsub.Pubsub
}

func New(store datastore.Datastore, pubsub pubsub.Pubsub) segmentrule.Management {
	store.EnforceDataType(segmentrule.Spec{})
	return &service{
		ruleStore: store,
		pubsub:    pubsub,
	}
}

func (s *service) Put(ctx context.Context, rule segmentrule.Spec) error {
	return s.ruleStore.RunInTransaction(ctx, func(ctx context.Context) error {
		original, exists, err := s.ruleStore.Get(ctx, rule.UID)
		if err != nil {
			return err
		}

		err = s.ruleStore.Put(ctx, rule.UID, rule)
		if err != nil {
			return err
		}

		if !exists {
			return s.pubsub.Publish(ctx, segmentrule.ManagementTopicName, segmentrule.CreatedEvent{RuleState: rule})
		} else {
			return s.pubsub.Publish(ctx, segmentrule.ManagementTopicName, segmentrule.ModifiedEvent{
				OldRuleState: original.(segmentrule.Spec),
				NewRuleState: rule,
			})
		}

		return nil
	})
}

func (s *service) Get(ctx context.Context, ruleUID string) (segmentrule.Spec, bool, error) {
	rule := segmentrule.Spec{}
	ruleExists := false
	var err error

	err = s.ruleStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.ruleStore.Get(ctx, ruleUID)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", ruleUID, err)
		}

		ruleExists = exists
		if !exists {
			return nil
		}

		rule = item.(segmentrule.Spec)

		return nil
	})
	if err != nil {
		return rule, false, err
	}

	return rule, ruleExists, nil
}

func (s *service) Remove(ctx context.Context, ruleUID string) error {
	return s.ruleStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.ruleStore.Get(ctx, ruleUID)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", ruleUID, err)
		}

		if exists {
			err = s.ruleStore.Remove(ctx, ruleUID)
			if err != nil {
				return fmt.Errorf("Error removing rule with uid %s: %s", ruleUID, err)
			}

			err = s.pubsub.Publish(ctx, segmentrule.ManagementTopicName, segmentrule.RemovedEvent{
				SegmentState: item.(segmentrule.Spec),
			})
			if err != nil {
				return fmt.Errorf("Error publishing RemovedEvent for rule %s: %s", ruleUID, err)
			}
		}
		return nil
	})
}

func (s *service) List(ctx context.Context) ([]segmentrule.Spec, error) {
	rules := []segmentrule.Spec{}
	var err error

	err = s.ruleStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.ruleStore.GetAll(ctx)
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

func (m *service) Preprov(ctx context.Context) error {
	err := m.Put(ctx, predefinedrules.YoungAgeSmsRule)
	if err != nil {
		return err
	}

	err = m.Put(ctx, predefinedrules.OldAgeEmailRule)
	if err != nil {
		return err
	}
	return nil
}
