package segmentmanagement

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type segmentManagement struct {
	segmentStore datastore.Datastore
	pubsub       pubsub.Pubsub
}

type SegmentManagement interface {
	segment.Management
}

func New(datastore datastore.Datastore, pubsub pubsub.Pubsub) SegmentManagement {
	return &segmentManagement{
		segmentStore: datastore,
		pubsub:       pubsub,
	}
}

func (s *segmentManagement) Put(ctx context.Context, segm segment.SegmentSpec) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		original, exists, err := s.segmentStore.Get(ctx, segm.UID)
		if err != nil {
			return err
		}

		err = s.segmentStore.Put(ctx, segm.UID, segm)
		if err != nil {
			return err
		}

		if !exists {
			return s.pubsub.Publish(ctx, segment.ManagementTopicName, segment.CreatedEvent{SegmentState: segm})
		} else {
			return s.pubsub.Publish(ctx, segment.ManagementTopicName, segment.ModifiedEvent{
				OldSegmentState: original.(segment.SegmentSpec),
				NewSegmentState: segm,
			})
		}

		return nil
	})
}

func (s *segmentManagement) Remove(ctx context.Context, segmentUID string) error {
	return s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		original, exists, err := s.segmentStore.Get(ctx, segmentUID)
		if err != nil {
			return err
		}

		if exists {
			err = s.segmentStore.Remove(ctx, segmentUID)
			if err != nil {
				return err
			}

			err = s.pubsub.Publish(ctx, segment.ManagementTopicName, segment.RemovedEvent{
				SegmentState: original.(segment.SegmentSpec),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *segmentManagement) Get(ctx context.Context, segmentUID string) (segment.SegmentSpec, bool, error) {
	var segm segment.SegmentSpec
	segmentExists := false
	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.segmentStore.Get(ctx, segmentUID)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("SegmentSpec with uid %s does not exist", segmentUID)
		}
		segm = item.(segment.SegmentSpec)
		segmentExists = exists

		return nil
	})
	if err != nil {
		return segm, false, err
	}
	return segm, segmentExists, nil
}

func (s *segmentManagement) List(ctx context.Context) ([]segment.SegmentSpec, error) {
	segments := []segment.SegmentSpec{}

	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, i := range items {
			segments = append(segments, i.(segment.SegmentSpec))
		}

		return nil
	})
	if err != nil {
		return segments, err
	}
	return segments, nil
}