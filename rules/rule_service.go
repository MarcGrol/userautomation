package rules

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"reflect"

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

func (s *userSegmentRuleService) Put(ctx context.Context, rule UserSegmentRule) error {
	return s.datastore.RunInTransaction(ctx, func(ctx context.Context) error {
		originalRule, exists, err := s.datastore.Get(ctx, rule.Name)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", rule.Name, err)
		}

		err = s.datastore.Put(ctx, rule.Name, rule)
		if err != nil {
			return fmt.Errorf("Error storing rule with uid %s: %s", rule.Name, err)
		}

		if !exists {
			err = s.pubsub.Publish(ctx, RuleTopicName, RuleCreatedEvent{
				State: rule,
			})
			if err != nil {
				return fmt.Errorf("Error publishing rule-created-event uid %s: %s", rule.Name, err)
			}

		} else if !reflect.DeepEqual(originalRule, rule) {
			err := s.pubsub.Publish(ctx, RuleTopicName, RuleModifiedEvent{
				OldState: originalRule.(UserSegmentRule),
				NewState: rule,
			})
			if err != nil {
				return fmt.Errorf("Error publishing rule-modified-event for user %s: %s", rule.Name, err)
			}
		} else {
			// rule unchanged: do not notify
		}

		return nil
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
		rule, exists, err := s.datastore.Get(ctx, ruleName)
		if err != nil {
			return fmt.Errorf("Error fetching rule with uid %s: %s", ruleName, err)
		}
		if exists {
			err = s.datastore.Remove(ctx, ruleName)
			if err != nil {
				return fmt.Errorf("Error deleting rule with uid %s: %s", ruleName, err)
			}

			err = s.pubsub.Publish(ctx, RuleTopicName, RuleRemovedEvent{
				State: rule.(UserSegmentRule),
			})
			if err != nil {
				return fmt.Errorf("Error publishing rule-created-event uid %s: %s", ruleName, err)
			}
		}
		return nil
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
