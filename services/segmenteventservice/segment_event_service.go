package segmenteventservice

import (
	"context"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/user"

	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type segmentEventHandler struct {
	pubsub      pubsub.Pubsub
	ruleService rule.SegmentRuleService
}

type SegmentEventService interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	segment.EventHandler
}

func NewSegmentEventService(pubsub pubsub.Pubsub, ruleService rule.SegmentRuleService) SegmentEventService {
	return &segmentEventHandler{
		pubsub:      pubsub,
		ruleService: ruleService,
	}
}

func (s *segmentEventHandler) IamSubscribing() {}

func (s *segmentEventHandler) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, segment.TopicName, s.OnEvent)
}

func (s *segmentEventHandler) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return segment.DispatchEvent(ctx, s, topic, event)
}

func (s *segmentEventHandler) OnUserAddedToSegment(ctx context.Context, event segment.UserAddedToSegmentEvent) error {
	return s.performUserActionForAllMatchingRules(ctx, event)
}

func (s *segmentEventHandler) OnUserRemovedFromSegment(ctx context.Context, event segment.UserRemovedFromSegmentEvent) error {
	return nil
}

func (s *segmentEventHandler) performUserActionForAllMatchingRules(ctx context.Context, event segment.UserAddedToSegmentEvent) error {
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return err
	}

	for _, r := range rules {
		if r.UserSegment.UID == event.SegmentUID {
			s.performUserAction(ctx, r, event.User)
		}
	}
	return nil
}

func (s *segmentEventHandler) performUserAction(ctx context.Context, r rule.UserSegmentRule, u user.User) error {

	// double check
	applicable, err := r.UserSegment.IsApplicableForUser(ctx, u)
	if err != nil {
		return err
	}
	if applicable {
		err = r.Action.Perform(ctx, action.UserAction{
			RuleUID:     r.UID,
			TriggerType: action.UserAddedToSegment,
			NewState:    &u,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
