package segmentchangeevaluator

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/infra/taskqueue"
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
	taskqueue   taskqueue.TaskQueue
	ruleService rule.RuleService
}

func New(pubsub pubsub.Pubsub, ruleService rule.RuleService, taskqueue taskqueue.TaskQueue) SegmentChangeEvaluator {
	return &segmentChangeEvaluator{
		pubsub:      pubsub,
		ruleService: ruleService,
		taskqueue:   taskqueue,
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
			act := action.UserAction{
				RuleUID:  r.UID,
				Reason:   0, // TODO
				OldState: nil,
				NewState: &event.User,
			}
			payload, err := json.MarshalIndent(act, "", "\t")
			if err != nil {
				return err
			}
			t := taskqueue.Task{
				Method:  "POST",
				URL:     "/api/action",
				Payload: string(payload),
			}
			return s.taskqueue.Enqueue(ctx, t)
		}
	}

	return nil
}
func (s *segmentChangeEvaluator) OnUserRemovedFromSegment(ctx context.Context, event segment.UserRemovedFromSegmentEvent) error {
	return fmt.Errorf("not implemented")
}
