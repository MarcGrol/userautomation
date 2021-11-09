package usercchangevaluator

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"

	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type UserChangeEvaluator interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	user.EventHandler
}

type userChangeEvaluator struct {
	pubsub      pubsub.Pubsub
	ruleService rule.RuleService
}

func New(pubsub pubsub.Pubsub, ruleService rule.RuleService) UserChangeEvaluator {
	return &userChangeEvaluator{
		pubsub:      pubsub,
		ruleService: ruleService,
	}
}

func (s *userChangeEvaluator) IamSubscribing() {}

func (s *userChangeEvaluator) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, user.ManagementTopicName, s.OnEvent)
}

func (s *userChangeEvaluator) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return user.DispatchEvent(ctx, s, topic, event)
}

func (s *userChangeEvaluator) OnUserCreated(ctx context.Context, event user.CreatedEvent) error {
	u := event.UserState

	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, r := range ruleSlice {
		executed, err := executeRuleForUser(ctx, r, u, action.ReasonIsUserAddedToSegment)
		if err != nil {
			return fmt.Errorf("Error executing r %s for user %s: %s", r.UID, u.UID, err)
		}
		if executed {
			s.onActionPerformed(ctx, r, u)
		}
	}
	return nil
}

func (s *userChangeEvaluator) OnUserModified(ctx context.Context, event user.ModifiedEvent) error {
	oldState := event.OldUserState
	newState := event.NewUserState

	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range ruleSlice {
		ruleApplicableBefore, err := rule.SegmentSpec.IsApplicableForUser(ctx, oldState)
		if err != nil {
			return fmt.Errorf("Error determining if rule %s is applicable for user %s: %s", rule.UID, newState.UID, err)
		}

		ruleApplicableAfter, err := rule.SegmentSpec.IsApplicableForUser(ctx, newState)
		if err != nil {
			return fmt.Errorf("Error determining if rule %s is applicable for user %s: %s", rule.UID, newState.UID, err)
		}

		if !ruleApplicableBefore && ruleApplicableAfter {
			err = rule.Action.Perform(ctx, action.UserAction{
				RuleUID:  rule.UID,
				Reason:   action.ReasonIsUserAddedToSegment,
				OldState: &oldState,
				NewState: &newState,
			})
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and user %s: %s", rule.UID, newState.UID, err)
			}
			s.onActionPerformed(ctx, rule, newState)

		} else {
			// do not execute the action if the user already belongs to this segment
		}
	}

	return nil
}

func (s *userChangeEvaluator) OnUserRemoved(ctx context.Context, event user.RemovedEvent) error {
	//
	return nil
}

func (s *userChangeEvaluator) onActionPerformed(ctx context.Context, rule rule.RuleSpec, user user.User) {
	// Should we keep track that this rule has fired for this user?
	// To prevent event being fired again when user re-enters again within particular time interval?
}

func executeRuleForUser(ctx context.Context, r rule.RuleSpec, user user.User, triggerType action.ReasonForAction) (bool, error) {
	applicable, err := r.SegmentSpec.IsApplicableForUser(ctx, user)
	if err != nil {
		return false, fmt.Errorf("Error determining if rule %s is applicable for u %s: %s", r.UID, user.UID, err)
	}

	if !applicable {
		return false, nil
	}

	err = r.Action.Perform(ctx, action.UserAction{
		RuleUID:  r.UID,
		Reason:   triggerType,
		OldState: nil,
		NewState: &user,
	})
	if err != nil {
		return false, fmt.Errorf("Error performing action for rule %s and useer %s: %s", r.UID, user.UID, err)
	}

	return true, nil
}
