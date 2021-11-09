package segmentchangeactioner

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/useraction"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/MarcGrol/userautomation/infra/taskqueue"
)

type SegmentChangeHandler interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	segment.UserEventHandler
}

type segmentChangeActioner struct {
	segment.UserEventHandler
	pubsub      pubsub.Pubsub
	taskqueue   taskqueue.TaskQueue
	ruleService rule.RuleService
}

func New(pubsub pubsub.Pubsub, ruleService rule.RuleService, taskqueue taskqueue.TaskQueue) SegmentChangeHandler {
	return &segmentChangeActioner{
		pubsub:      pubsub,
		ruleService: ruleService,
		taskqueue:   taskqueue,
	}
}

func (s *segmentChangeActioner) IamSubscribing() {}

func (s *segmentChangeActioner) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, segment.UserTopicName, s.OnEvent)
}

func (s *segmentChangeActioner) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return segment.DispatchUserEvent(ctx, s, topic, event)
}

func (s *segmentChangeActioner) OnUserAddedToSegment(ctx context.Context, event segment.UserAddedToSegmentEvent) error {
	// find actions related to this segment
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return err
	}

	for _, r := range rules {
		if r.SegmentSpec.UID == event.SegmentUID {
			act := useraction.UserAction{
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
func (s *segmentChangeActioner) OnUserRemovedFromSegment(ctx context.Context, event segment.UserRemovedFromSegmentEvent) error {
	return fmt.Errorf("not implemented")
}
