package segmentmanagement

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/coredata/predefinedsegments"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestSegmentManagement(t *testing.T) {
	ctx := context.TODO()

	t.Run("create segment", func(t *testing.T) {
		// setup

		segmentStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, ps)

		// given
		nothing()

		// when
		err := sut.Put(ctx, initialSegment())

		// then
		assert.NoError(t, err)
		assert.Len(t, listSegment(ctx, t, sut), 1)
		assert.Equal(t, "young users segment", getSegment(ctx, t, sut).Description)
		assert.Equal(t, initialSegment().UID, ps.Publications[0].Event.(segment.CreatedEvent).SegmentState.UID)
	})

	t.Run("modify segment", func(t *testing.T) {
		// setup
		segmentStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, ps)

		// given
		sut.Put(ctx, initialSegment())

		// when
		err := sut.Put(ctx, modifiedSegment())

		// then
		assert.NoError(t, err)
		assert.Len(t, listSegment(ctx, t, sut), 1)
		assert.Equal(t, modifiedSegment().Description, getSegment(ctx, t, sut).Description)
		assert.Equal(t, initialSegment().Description, ps.Publications[1].Event.(segment.ModifiedEvent).OldSegmentState.Description)
		assert.Equal(t, modifiedSegment().Description, ps.Publications[1].Event.(segment.ModifiedEvent).NewSegmentState.Description)
	})

	t.Run("remove segment", func(t *testing.T) {
		// setup
		segmentStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, ps)

		// given
		sut.Put(ctx, initialSegment())

		// when
		err := sut.Remove(ctx, initialSegment().UID)

		// then
		assert.NoError(t, err)
		assert.False(t, existSegment(ctx, t, sut))
		assert.Len(t, listSegment(ctx, t, sut), 0)
		assert.Equal(t, initialSegment().UID, ps.Publications[1].Event.(segment.RemovedEvent).SegmentState.UID)
	})
}

func setupMocks(t *testing.T) (*datastore.DatastoreStub, *pubsub.PubsubStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	storeSub := datastore.NewDatastoreStub()
	ps := pubsub.NewPubsubStub()
	return storeSub, ps, ctrl
}

func initialSegment() segment.Spec {
	return predefinedsegments.YoungAgeSegment
}

func modifiedSegment() segment.Spec {
	segm := predefinedsegments.OldAgeSegment
	segm.UID = predefinedsegments.YoungAgeSegment.UID
	return segm
}

func existSegment(ctx context.Context, t *testing.T, sut segment.Management) bool {
	_, exists, err := sut.Get(ctx, initialSegment().UID)
	if err != nil {
		t.Error(err)
	}
	return exists
}

func getSegment(ctx context.Context, t *testing.T, sut segment.Management) segment.Spec {
	segm, exists, err := sut.Get(ctx, initialSegment().UID)
	if err != nil || !exists {
		t.Error(err)
	}
	return segm
}

func listSegment(ctx context.Context, t *testing.T, sut segment.Management) []segment.Spec {
	segments, err := sut.List(ctx)
	if err != nil {
		t.Error(err)
	}
	return segments
}

func nothing() {}
