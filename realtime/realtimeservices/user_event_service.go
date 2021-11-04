package realtimeservices

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

type userEventHandler struct {
	pubsub      realtimecore.Pubsub
	ruleService realtimecore.SegmentRuleService
}

func NewUserEventService(pubsub realtimecore.Pubsub, ruleService realtimecore.SegmentRuleService) realtimecore.UserEventService {
	return &userEventHandler{
		pubsub:      pubsub,
		ruleService: ruleService,
	}
}

func (s *userEventHandler) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, "user", s.OnEvent)
}

func (s *userEventHandler) OnEvent(ctx context.Context, topic string, event interface{}) error {
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	switch e := event.(type) {
	case realtimecore.UserCreatedEvent:
		return s.onUserCreated(ctx, rules, e.State)
	case realtimecore.UserModifiedEvent:
		return s.onUserModified(ctx, rules, e.OldState, e.NewState)
	case realtimecore.UserRemovedEvent:
		return s.onUserRemoved(ctx, rules, e.State)
	default:
		return fmt.Errorf("Event %+v not supported", e)
	}
}

func (s *userEventHandler) onUserCreated(ctx context.Context, rules []realtimecore.UserSegmentRule, user realtimecore.User) error {
	for _, rule := range rules {
		applicable, err := rule.IsApplicableForUser(ctx, user)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for user: %s", err)
		}
		if applicable {
			err = rule.PerformAction(ctx, rule.Name, realtimecore.UserCreated, nil, &user)
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and useer %s: %s", rule.Name, user.UID, err)
			}
			s.onActionPerformed(ctx, rule, user)
		}
	}
	return nil
}

func (s *userEventHandler) onUserModified(ctx context.Context, rules []realtimecore.UserSegmentRule, oldState realtimecore.User, newState realtimecore.User) error {
	for _, rule := range rules {
		ruleApplicableBefore, err := rule.IsApplicableForUser(ctx, oldState)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for newState: %s", err)
		}

		ruleApplicableAfter, err := rule.IsApplicableForUser(ctx, newState)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for newState: %s", err)
		}

		if !ruleApplicableBefore && ruleApplicableAfter {
			err = rule.PerformAction(ctx, rule.Name, realtimecore.UserModified, &oldState, &newState)
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and us	er %s: %s", rule.Name, newState.UID, err)
			}
			s.onActionPerformed(ctx, rule, newState)

		} else {
			// do not execute the action if the user already belongs to this segment
		}
	}

	return nil
}

func (s *userEventHandler) onUserRemoved(ctx context.Context, rules []realtimecore.UserSegmentRule, user realtimecore.User) error {
	for _, rule := range rules {
		applicable, err := rule.IsApplicableForUser(ctx, user)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for user: %s", err)
		}
		if applicable {
			err = rule.PerformAction(ctx, rule.Name, realtimecore.UserRemoved, &user, nil)
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and useer %s: %s", rule.Name, user.UID, err)
			}
			s.onActionPerformed(ctx, rule, user)
		}
	}
	return nil
}

func (s *userEventHandler) onActionPerformed(ctx context.Context, rule realtimecore.UserSegmentRule, user realtimecore.User) {
	// Should we keep track that this rule has fired for this user?
	// To prevent event being dfire again when user re-enters again within particular time interval?
}