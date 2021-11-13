package segmentmanagement

import (
	"context"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/coredata/predefinedsegments"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
)

type service struct {
	segmentStore datastore.Datastore
	pubsub       pubsub.Pubsub
}

func New(store datastore.Datastore, pubsub pubsub.Pubsub) segment.Management {
	store.EnforceDataType(segment.Spec{})
	return &service{
		segmentStore: store,
		pubsub:       pubsub,
	}
}

func (s *service) Put(ctx context.Context, segm segment.Spec) error {
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
				OldSegmentState: original.(segment.Spec),
				NewSegmentState: segm,
			})
		}

		return nil
	})
}

func (s *service) Remove(ctx context.Context, segmentUID string) error {
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
				SegmentState: original.(segment.Spec),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *service) Get(ctx context.Context, segmentUID string) (segment.Spec, bool, error) {
	var segm segment.Spec
	segmentExists := false
	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		item, exists, err := s.segmentStore.Get(ctx, segmentUID)
		if err != nil {
			return err
		}

		segmentExists = exists
		if !exists {
			return nil
		}

		segm = item.(segment.Spec)

		return nil
	})
	if err != nil {
		return segm, false, err
	}
	return segm, segmentExists, nil
}

func (s *service) List(ctx context.Context) ([]segment.Spec, error) {
	segments := []segment.Spec{}

	err := s.segmentStore.RunInTransaction(ctx, func(ctx context.Context) error {
		items, err := s.segmentStore.GetAll(ctx)
		if err != nil {
			return err
		}

		for _, i := range items {
			segments = append(segments, i.(segment.Spec))
		}

		return nil
	})
	if err != nil {
		return segments, err
	}
	return segments, nil
}

func (m *service) Preprov(ctx context.Context) error {
	err := m.Put(ctx, predefinedsegments.YoungAgeSegment)
	if err != nil {
		return err
	}

	err = m.Put(ctx, predefinedsegments.OldAgeSegment)
	if err != nil {
		return err
	}

	err = m.Put(ctx, predefinedsegments.FirstNameStartsWithMSegment)
	if err != nil {
		return err
	}

	return nil
}
