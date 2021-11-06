package segmentservice

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"

	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type segmentService struct {
	datastore datastore.Datastore
	pubsub    pubsub.Pubsub
}

type SegmentService interface {
	// Flags that this service is an event consumer
	pubsub.SubscribingService
	// Early warning system. This service will break when "users"-module introduces new events.
	// In this case this service should also introduce these new events.
	user.UserEventHandler
	segment.UserSegmentService
}

func NewSegmentService(datastore datastore.Datastore, pubsub pubsub.Pubsub) SegmentService {
	return &segmentService{
		datastore: datastore,
		pubsub:    pubsub,
	}
}

func (s *segmentService) IamSubscribing() {}

func (s *segmentService) Subscribe(ctx context.Context) error {
	return s.pubsub.Subscribe(ctx, user.UserTopicName, s.OnEvent)
}

func (s *segmentService) OnEvent(ctx context.Context, topic string, event interface{}) error {
	return user.DispatchEvent(ctx, s, topic, event)
}

func (s *segmentService) OnUserCreated(ctx context.Context, user user.User) error {
	// TODO check if user must be added to segments

	return fmt.Errorf("Not implemented")
}

func (s *segmentService) OnUserModified(ctx context.Context, oldState user.User, newState user.User) error {
	// TODO check if user must be added to or removed from segments

	return fmt.Errorf("Not implemented")
}

func (s *segmentService) OnUserRemoved(ctx context.Context, user user.User) error {
	// TODO check if user must be removed from segments

	return fmt.Errorf("Not implemented")
}

func (s *segmentService) Put(ctx context.Context, userSegment segment.UserSegmentDefinition) error {
	// TODO re-evaluate all users that belong to this segment

	// TODO use datastore to persist

	return fmt.Errorf("Not implemented")
}

func (s *segmentService) Get(ctx context.Context, userSegmentUID string) (segment.UserSegmentDefinition, bool, error) {
	// TODO use datastore to persist
	return segment.UserSegmentDefinition{}, false, fmt.Errorf("Not implemented")
}

func (s *segmentService) List(ctx context.Context) ([]segment.UserSegmentDefinition, error) {
	// TODO use datastore to replace
	return []segment.UserSegmentDefinition{}, fmt.Errorf("Not implemented")
}

func (s *segmentService) Remove(ctx context.Context, userUID string) error {
	// TODO use datastore to remove

	return fmt.Errorf("Not implemented")
}

func (s *segmentService) GetUsersForSegment(ctx context.Context, userSegmentUID string) ([]user.User, error) {
	return []user.User{}, fmt.Errorf("Not implemented")
}
