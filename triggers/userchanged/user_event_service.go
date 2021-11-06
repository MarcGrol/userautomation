package userchanged

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/action"

	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/users"
)

type userEventHandler struct {
	pubsub      pubsub.Pubsub
	ruleService rules.SegmentRuleService
}

type UserEventService interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	users.UserEventHandler
}

func NewUserEventService(pubsub pubsub.Pubsub, ruleService rules.SegmentRuleService) UserEventService {
	return &userEventHandler{
		pubsub:      pubsub,
		ruleService: ruleService,
	}
}

func (s *userEventHandler) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, users.UserTopicName, s.OnEvent)
}

func (s *userEventHandler) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return users.DispatchEvent(ctx, s, topic, event)
}

func (s *userEventHandler) OnUserCreated(ctx context.Context, user users.User) error {
	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range ruleSlice {
		applicable, err := rule.UserSegment.IsApplicableForUser(ctx, user)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for user: %s", err)
		}
		if applicable {
			err = rule.Action.Perform(ctx, action.UserAction{
				RuleName:       rule.Name,
				UserChangeType: action.UserCreated,
				OldState:       nil,
				NewState:       &user,
			})
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and user %s: %s", rule.Name, user.UID, err)
			}
			s.onActionPerformed(ctx, rule, user)
		}
	}
	return nil
}

func (s *userEventHandler) OnUserModified(ctx context.Context, oldState users.User, newState users.User) error {
	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range ruleSlice {
		ruleApplicableBefore, err := rule.UserSegment.IsApplicableForUser(ctx, oldState)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for newState: %s", err)
		}

		ruleApplicableAfter, err := rule.UserSegment.IsApplicableForUser(ctx, newState)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for newState: %s", err)
		}

		if !ruleApplicableBefore && ruleApplicableAfter {
			err = rule.Action.Perform(ctx, action.UserAction{
				RuleName:       rule.Name,
				UserChangeType: action.UserModified,
				OldState:       &oldState,
				NewState:       &newState,
			})
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and userService	er %s: %s", rule.Name, newState.UID, err)
			}
			s.onActionPerformed(ctx, rule, newState)

		} else {
			// do not execute the action if the user already belongs to this segment
		}
	}

	return nil
}

func (s *userEventHandler) OnUserRemoved(ctx context.Context, user users.User) error {
	ruleSlice, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range ruleSlice {
		applicable, err := rule.UserSegment.IsApplicableForUser(ctx, user)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for user: %s", err)
		}
		if applicable {
			err = rule.Action.Perform(ctx, action.UserAction{
				RuleName:       rule.Name,
				UserChangeType: action.UserRemoved,
				OldState:       nil,
				NewState:       &user,
			})
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and useer %s: %s", rule.Name, user.UID, err)
			}
			s.onActionPerformed(ctx, rule, user)
		}
	}
	return nil
}

func (s *userEventHandler) onActionPerformed(ctx context.Context, rule rules.UserSegmentRule, user users.User) {
	// Should we keep track that this rule has fired for this user?
	// To prevent event being dfire again when user re-enters again within particular time interval?
}
