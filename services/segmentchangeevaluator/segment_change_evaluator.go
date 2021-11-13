package segmentchangeevaluator

import (
	"context"
	"fmt"
	"log"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type Service interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	segment.UserEventHandler
}

type service struct {
	segment.UserEventHandler
	pubsub      pubsub.Pubsub
	ruleService segmentrule.Service
}

func New(pubsub pubsub.Pubsub, ruleService segmentrule.Service) Service {
	return &service{
		pubsub:      pubsub,
		ruleService: ruleService,
	}
}

func (s *service) IamSubscribing() {}

func (s *service) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, segment.UserTopicName, s.OnEvent)
}

func (s *service) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return segment.DispatchUserEvent(ctx, s, topic, event)
}

func (s *service) OnUserAddedToSegment(ctx context.Context, event segment.UserAddedToSegmentEvent) error {
	// find actions related to this segment
	rules, err := s.ruleService.List(ctx)
	if err != nil {
		return err
	}

	for _, r := range rules {
		if r.SegmentSpec.UID == event.SegmentUID {

			if !event.User.HasAttributes(r.ActionSpec.MandatoryUserAttributes) {
				log.Printf("User %s is missing madatory attributes for action %s", event.User.UID, r.ActionSpec.Name)
				continue
			}

			err := s.pubsub.Publish(ctx, usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
				Task: usertask.Spec{
					UID:        "", // TODO identify each triggered rule uninquely
					RuleUID:    r.UID,
					ActionSpec: r.ActionSpec,
					Reason:     usertask.ReasonUserAddedToSegment,
					User:       event.User,
				},
			})
			if err != nil {
				return fmt.Errorf("Error publishing user-task for rule %s and user %s: %s", r.UID, event.User.UID, err)
			}
		}
	}

	return nil
}

func (s *service) OnUserRemovedFromSegment(ctx context.Context, event segment.UserRemovedFromSegmentEvent) error {
	return fmt.Errorf("not implemented")
}
