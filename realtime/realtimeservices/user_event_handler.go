package realtimeservices

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

type userEventHandler struct {
	ruleService realtimecore.SegmentRuleService
}

func NewUserEventHandler(ruleService realtimecore.SegmentRuleService) realtimecore.UserEventService {
	return &userEventHandler{
		ruleService: ruleService,
	}
}

func (s *userEventHandler) OnUserCreated(ctx context.Context, user realtimecore.User) error {
	return s.onUserEvent(ctx, realtimecore.UserCreated, user)
}

func (s *userEventHandler) OnUserModified(ctx context.Context, before realtimecore.User, user realtimecore.User) error {
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range rules {
		applicableBefore, err := rule.IsApplicableForUser(ctx, before)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for user: %s", err)
		}
		applicableAfter, err := rule.IsApplicableForUser(ctx, user)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for user: %s", err)
		}
		if !applicableBefore && applicableAfter {
			err = rule.PerformAction(ctx, rule.Name, realtimecore.UserModified, user)
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and us	er %s: %s", rule.Name, user.FullName, err)
			}
		} else {
			// do not execute the action if the user already belongs to this segment
		}
	}

	return nil
}

func (s *userEventHandler) OnUserRemoved(c context.Context, user realtimecore.User) error {
	return s.onUserEvent(c, realtimecore.UserRemoved, user)
}

func (s *userEventHandler) onUserEvent(ctx context.Context, status realtimecore.UserStatus, user realtimecore.User) error {
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

	for _, rule := range rules {
		applicable, err := rule.IsApplicableForUser(ctx, user)
		if err != nil {
			return fmt.Errorf("Error determining if rule is applicable for user: %s", err)
		}
		if applicable {
			err = rule.PerformAction(ctx, rule.Name, status, user)
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and useer %s: %s", rule.Name, user.FullName, err)
			}
		}
	}
	return nil
}
