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
	user.EventHandler
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

	for _, r := range ruleSlice {
		executed, err := executeRuleForUser(ctx, r, u, action.UserCreated)
		if err != nil {
			return fmt.Errorf("Error executing r %s for user %s: %s", r.UID, u.UID, err)
		}
		if executed {
			s.onActionPerformed(ctx, r, u)
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
				RuleName:    rule.UID,
				TriggerType: action.UserModified,
				OldState:    &oldState,
				NewState:    &newState,
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

	for _, r := range ruleSlice {
		executed, err := executeRuleForUser(ctx, r, u, action.UserRemoved)
		if err != nil {
			return fmt.Errorf("Error executing r %s for user %s: %s", r.UID, u.UID, err)
		}
		if executed {
			s.onActionPerformed(ctx, r, u)
		}
	}
	return nil
}

func (s *userEventHandler) onActionPerformed(ctx context.Context, rule rule.UserSegmentRule, user user.User) {
	// Should we keep track that this rule has fired for this user?
	// To prevent event being fired again when user re-enters again within particular time interval?
}

func executeRuleForUser(ctx context.Context, r rule.UserSegmentRule, user user.User, triggerType action.TriggerType) (bool, error) {

	// TODO this breaks tests
	//if (r.TriggerKindMask & rule.TriggerUserChange) == 1 {
	//	return false, nil
	//}

	applicable, err := r.UserSegment.IsApplicableForUser(ctx, user)
	if err != nil {
		return false, fmt.Errorf("Error determining if rule %s is applicable for u %s: %s", r.UID, user.UID, err)
	}

	if !applicable {
		return false, nil
	}

	err = r.Action.Perform(ctx, action.UserAction{
		RuleName:    r.UID,
		TriggerType: triggerType,
		OldState:    nil,
		NewState:    &user,
	})
	if err != nil {
		return false, fmt.Errorf("Error performing action for rule %s and useer %s: %s", r.UID, user.UID, err)
	}

	return true, nil
}
