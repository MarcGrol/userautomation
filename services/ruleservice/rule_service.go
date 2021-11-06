package ruleservice

import (
	"context"
	"fmt"
	rule2 "github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"reflect"

	"github.com/MarcGrol/userautomation/infra/datastore"
)

type userSegmentRuleService struct {
	datastore datastore.Datastore
	pubsub    pubsub.Pubsub
}

func NewUserSegmentRuleService(datastore datastore.Datastore, pubsub pubsub.Pubsub) rule2.SegmentRuleService {
	return &userSegmentRuleService{
		datastore: datastore,
		pubsub:    pubsub, // TODO might signal rule change to dedicated service
	}
}

func (s *userSegmentRuleService) Put(ctx context.Context, rule rule2.UserSegmentRule) error {
	return s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		originalRule, exists, err := s.datastore.Get(ctx, rule.UID)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", rule.UID, err)
		}

		err = s.datastore.Put(ctx, rule.UID, rule)
		if err != nil {
			return fmt.Errorf("Error storing rule with uid %s: %s", rule.UID, err)
		}

		if !exists {
			err = s.pubsub.Publish(ctx, rule2.RuleTopicName, rule2.RuleCreatedEvent{
				State: rule,
			})
			if err != nil {
				return fmt.Errorf("Error publishing rule-created-event uid %s: %s", rule.UID, err)
			}

		} else if !reflect.DeepEqual(originalRule, rule) {
			err := s.pubsub.Publish(ctx, rule2.RuleTopicName, rule2.RuleModifiedEvent{
				OldState: originalRule.(rule2.UserSegmentRule),
				NewState: rule,
			})
			if err != nil {
				return fmt.Errorf("Error publishing rule-modified-event for user %s: %s", rule.UID, err)
			}
		} else {
			// rule unchanged: do not notify
		}

		return nil
	})
}

func (s *userSegmentRuleService) Get(ctx context.Context, ruleUID string) (rule2.UserSegmentRule, bool, error) {
	rule := rule2.UserSegmentRule{}
	ruleExists := false
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.datastore.Get(ctx, ruleUID)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", ruleUID, err)
		}
		rule = item.(rule2.UserSegmentRule)
		ruleExists = exists

		return nil
	})
	if err != nil {
		return rule, false, err
	}

	return rule, ruleExists, nil
}

func (s *userSegmentRuleService) Delete(ctx context.Context, ruleUID string) error {
	return s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		rule, exists, err := s.datastore.Get(ctx, ruleUID)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", ruleUID, err)
		}
		if exists {
			err = s.datastore.Remove(ctx, ruleUID)
			if err != nil {
				return fmt.Errorf("Error deleting rule with uid %s: %s", ruleUID, err)
			}

			err = s.pubsub.Publish(ctx, rule2.RuleTopicName, rule2.RuleRemovedEvent{
				State: rule.(rule2.UserSegmentRule),
			})
			if err != nil {
				return fmt.Errorf("Error publishing rule-created-event uid %s: %s", ruleUID, err)
			}
		}
		return nil
	})
}

func (s *userSegmentRuleService) List(ctx context.Context) ([]rule2.UserSegmentRule, error) {
	rules := []rule2.UserSegmentRule{}
	var err error

	err = s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.datastore.GetAll(ctx)
		if err != nil {
			return fmt.Errorf("Error fetching all rules: %s", err)
		}

		for _, item := range items {
			rules = append(rules, item.(rule2.UserSegmentRule))
		}
		return nil
	})
	if err != nil {
		return rules, err
	}

	return rules, nil
}
