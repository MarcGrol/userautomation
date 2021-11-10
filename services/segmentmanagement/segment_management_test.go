package segmentmanagement

import (
	"context"
	supportedrules "github.com/MarcGrol/userautomation/coredata/supportedsegments"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/MarcGrol/userautomation/core/segment"
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
		createSegment(ctx, t, sut)

		// then
		assert.Len(t, listSegment(ctx, t, sut), 1)
		assert.Equal(t, "young users segment", getSegment(ctx, t, sut).Description)
		assert.Equal(t, 1, ps.PublicationCount)
	})

	t.Run("modify segment", func(t *testing.T) {
		// setup
		segmentStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, ps)

		// given
		createSegment(ctx, t, sut)

		// when
		modifySegment(ctx, t, sut)

		// then
		assert.Len(t, listSegment(ctx, t, sut), 1)
		assert.Equal(t, "old users segment", getSegment(ctx, t, sut).Description)
		assert.Equal(t, 2, ps.PublicationCount)
	})

	t.Run("remove segment", func(t *testing.T) {
		// setup
		segmentStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, ps)

		// given
		createSegment(ctx, t, sut)

		// when
		removeSegment(ctx, t, sut)

		// then
		assert.Len(t, listSegment(ctx, t, sut), 0)
		assert.Equal(t, 2, ps.PublicationCount)
	})
}

func setupMocks(t *testing.T) (*datastore.DatastoreStub, *pubsub.PubsubStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	storeSub := datastore.NewDatastoreStub()
	ps := pubsub.NewPubsubStub()
	return storeSub, ps, ctrl
}

func initialSegment() segment.Spec {
	return supportedrules.YoungAgeSegment
}

func createSegment(ctx context.Context, t *testing.T, sut segment.Management) {
	err := sut.Put(ctx, initialSegment())
	if err != nil {
		t.Error(err)
	}
}

func modifiedSegment() segment.Spec {
	segm := supportedrules.OldAgeSegment
	segm.UID = supportedrules.YoungAgeSegment.UID
	return segm
}

func modifySegment(ctx context.Context, t *testing.T, sut segment.Management) {
	err := sut.Put(ctx, modifiedSegment())
	if err != nil {
		t.Error(err)
	}
}
func removeSegment(ctx context.Context, t *testing.T, sut segment.Management) {
	err := sut.Remove(ctx, initialSegment().UID)
	if err != nil {
		t.Error(err)
	}
}

func existsSegment(ctx context.Context, t *testing.T, sut segment.Management) bool {
	_, exists, err := sut.Get(ctx, initialSegment().UID)
	if err != nil || !exists {
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
