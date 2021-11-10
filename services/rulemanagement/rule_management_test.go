package rulemanagement

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/coredata/predefinedrules"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestRuleManagement(t *testing.T) {
	ctx := context.TODO()

	t.Run("create rule", func(t *testing.T) {
		// setup

		ruleStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(ruleStore, ps)

		// given
		nothing()

		// when
		err := sut.Put(ctx, initialRule())

		// then
		assert.NoError(t, err)
		assert.Len(t, listRule(ctx, t, sut), 1)
		assert.Equal(t, "Send sms to young users", getRule(ctx, t, sut).Description)
		assert.Equal(t, initialRule().UID, ps.Publications[0].Event.(segmentrule.CreatedEvent).RuleState.UID)
	})

	t.Run("modify rule", func(t *testing.T) {
		// setup
		ruleStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(ruleStore, ps)

		// given
		sut.Put(ctx, initialRule())

		// when
		err := sut.Put(ctx, modifiedRule())

		// then
		assert.NoError(t, err)
		assert.Len(t, listRule(ctx, t, sut), 1)
		assert.Equal(t, modifiedRule().Description, getRule(ctx, t, sut).Description)
		assert.Equal(t, initialRule().Description, ps.Publications[1].Event.(segmentrule.ModifiedEvent).OldRuleState.Description)
		assert.Equal(t, modifiedRule().Description, ps.Publications[1].Event.(segmentrule.ModifiedEvent).NewRuleState.Description)
	})

	t.Run("remove rule", func(t *testing.T) {
		// setup
		ruleStore, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(ruleStore, ps)

		// given
		sut.Put(ctx, initialRule())

		// when
		err := sut.Remove(ctx, initialRule().UID)

		// then
		assert.NoError(t, err)
		assert.False(t, existRule(ctx, t, sut))
		assert.Len(t, listRule(ctx, t, sut), 0)
		assert.Equal(t, initialRule().UID, ps.Publications[1].Event.(segmentrule.RemovedEvent).SegmentState.UID)
	})
}

func setupMocks(t *testing.T) (*datastore.DatastoreStub, *pubsub.PubsubStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	storeSub := datastore.NewDatastoreStub()
	ps := pubsub.NewPubsubStub()
	return storeSub, ps, ctrl
}

func initialRule() segmentrule.Spec {
	return predefinedrules.YoungAgeSmsRule
}

func modifiedRule() segmentrule.Spec {
	segm := predefinedrules.OldAgeEmailRule
	segm.UID = predefinedrules.YoungAgeSmsRule.UID
	return segm
}

func existRule(ctx context.Context, t *testing.T, sut segmentrule.Service) bool {
	_, exists, err := sut.Get(ctx, initialRule().UID)
	if err != nil {
		t.Error(err)
	}
	return exists
}

func getRule(ctx context.Context, t *testing.T, sut segmentrule.Service) segmentrule.Spec {
	segm, exists, err := sut.Get(ctx, initialRule().UID)
	if err != nil || !exists {
		t.Error(err)
	}
	return segm
}

func listRule(ctx context.Context, t *testing.T, sut segmentrule.Service) []segmentrule.Spec {
	rules, err := sut.List(ctx)
	if err != nil {
		t.Error(err)
	}
	return rules
}

func nothing() {}
