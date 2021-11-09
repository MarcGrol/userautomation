package segmentchangeevaluator

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type SegmentChangeEvaluator interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	segment.UserEventHandler
}

type segmentChangeEvaluator struct {
	segment.UserEventHandler
	pubsub      pubsub.Pubsub
	ruleService rule.RuleService
}

func New(pubsub pubsub.Pubsub, ruleService rule.RuleService) SegmentChangeEvaluator {
	return &segmentChangeEvaluator{
		pubsub:      pubsub,
		ruleService: ruleService,
	}
}

func (s *segmentChangeEvaluator) IamSubscribing() {}

func (s *segmentChangeEvaluator) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, segment.UserTopicName, s.OnEvent)
}

func (s *segmentChangeEvaluator) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return segment.DispatchUserEvent(ctx, s, topic, event)
}

func (s *segmentChangeEvaluator) OnUserAddedToSegment(ctx context.Context, event segment.UserAddedToSegmentEvent) error {
	// find actions related to this segment
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return err
	}

	for _, r := range rules {
		if r.SegmentSpec.UID == event.SegmentUID {
			err := s.pubsub.Publish(ctx, usertask.TopicName, usertask.UserTask{
				RuleSpec: r,
				Reason:   0,
				User:     event.User,
			})
			if err != nil {
				return fmt.Errorf("Error publishing user-task for rule %s and user %s: %s", r.UID, event.User.UID, err)
			}
		}
	}

	return nil
}

func (s *segmentChangeEvaluator) OnUserRemovedFromSegment(ctx context.Context, event segment.UserRemovedFromSegmentEvent) error {
	return fmt.Errorf("not implemented")
}
