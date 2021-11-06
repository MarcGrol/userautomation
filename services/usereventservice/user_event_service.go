package usereventservice

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"

	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type userEventHandler struct {
	pubsub      pubsub.Pubsub
	ruleService rule.SegmentRuleService
}

type UserEventService interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	user.UserEventHandler
}

func NewUserEventService(pubsub pubsub.Pubsub, ruleService rule.SegmentRuleService) UserEventService {
	return &userEventHandler{
		pubsub:      pubsub,
		ruleService: ruleService,
	}
}

func (s *userEventHandler) IamSubscribing() {}

func (s *userEventHandler) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, user.UserTopicName, s.OnEvent)
}

func (s *userEventHandler) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return user.DispatchEvent(ctx, s, topic, event)
}

func (s *userEventHandler) OnUserCreated(ctx context.Context, u user.User) error {
	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range ruleSlice {
		applicable, err := rule.UserSegment.IsApplicableForUser(ctx, u)
		if err != nil {
			return fmt.Errorf("Error determining if rule %s is applicable for u %s: %s", rule.UID, u.UID, err)
		}
		if applicable {
			err = rule.Action.Perform(ctx, action.UserAction{
				RuleName:       rule.UID,
				UserChangeType: action.UserCreated,
				OldState:       nil,
				NewState:       &u,
			})
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and u %s: %s", rule.UID, u.UID, err)
			}
			s.onActionPerformed(ctx, rule, u)
		}
	}
	return nil
}

func (s *userEventHandler) OnUserModified(ctx context.Context, oldState user.User, newState user.User) error {
	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range ruleSlice {
		ruleApplicableBefore, err := rule.UserSegment.IsApplicableForUser(ctx, oldState)
		if err != nil {
			return fmt.Errorf("Error determining if rule %s is applicable for user %s: %s", rule.UID, newState.UID, err)
		}

		ruleApplicableAfter, err := rule.UserSegment.IsApplicableForUser(ctx, newState)
		if err != nil {
			return fmt.Errorf("Error determining if rule %s is applicable for user %s: %s", rule.UID, newState.UID, err)
		}

		if !ruleApplicableBefore && ruleApplicableAfter {
			err = rule.Action.Perform(ctx, action.UserAction{
				RuleName:       rule.UID,
				UserChangeType: action.UserModified,
				OldState:       &oldState,
				NewState:       &newState,
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

func (s *userEventHandler) OnUserRemoved(ctx context.Context, u user.User) error {
	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range ruleSlice {
		applicable, err := rule.UserSegment.IsApplicableForUser(ctx, u)
		if err != nil {
			return fmt.Errorf("Error determining if rule %s is applicable for u %s: %s", rule.UID, u.UID, err)
		}
		if applicable {
			err = rule.Action.Perform(ctx, action.UserAction{
				RuleName:       rule.UID,
				UserChangeType: action.UserRemoved,
				OldState:       nil,
				NewState:       &u,
			})
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and useer %s: %s", rule.UID, u.UID, err)
			}
			s.onActionPerformed(ctx, rule, u)
		}
	}
	return nil
}

func (s *userEventHandler) onActionPerformed(ctx context.Context, rule rule.UserSegmentRule, user user.User) {
	// Should we keep track that this rule has fired for this user?
	// To prevent event being dfire again when user re-enters again within particular time interval?
}
