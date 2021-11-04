package realtimeservices

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

type userEventHandler struct {
	pubsub realtimecore.Pubsub
	ruleService realtimecore.SegmentRuleService
}

func NewUserEventHandler(pubsub realtimecore.Pubsub, ruleService realtimecore.SegmentRuleService) realtimecore.UserEventService {
	return &userEventHandler{
		pubsub: pubsub,
		ruleService: ruleService,
	}
}

func (s *userEventHandler) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, "user", s.OnEvent)
}

func (s *userEventHandler) OnEvent(ctx context.Context, topic string, event interface{}) error {
	switch e := event.(type) {
	case realtimecore.UserCreatedEvent:
		return s.onUserCreated(ctx, e.State)
	case realtimecore.UserModifiedEvent:
		return s.onUserModified(ctx, e.OldState, e.NewState)
	case realtimecore.UserRemovedEvent:
		return s.onUserRemoved(ctx, e.State)
	default:
		return fmt.Errorf("Event %+v not supported", e)
	}
}

func (s *userEventHandler) onUserCreated(ctx context.Context, user realtimecore.User) error {
	return s.onUserEvent(ctx, realtimecore.UserCreated, user)
}

func (s *userEventHandler) onUserModified(ctx context.Context, oldState realtimecore.User, newState realtimecore.User) error {
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return fmt.Errorf("Error fetching rules: %s", err)
	}

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
			err = rule.PerformAction(ctx, rule.Name, realtimecore.UserModified, newState)
			if err != nil {
				return fmt.Errorf("Error performing action for rule %s and us	er %s: %s", rule.Name, newState.UID, err)
			}
		} else {
			// do not execute the action if the user already belongs to this segment
		}
	}

	return nil
}

func (s *userEventHandler) onUserRemoved(c context.Context, user realtimecore.User) error {
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
				return fmt.Errorf("Error performing action for rule %s and useer %s: %s", rule.Name, user.UID	, err)
			}
			// Should we keep track that this rule has fired for this user?
			// To prevent event being dfire again when user re-enters again within particular time interval?
		}
	}
	return nil
}
