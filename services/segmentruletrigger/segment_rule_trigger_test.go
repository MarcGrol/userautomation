package segmentruletrigger

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/coredata/predefinedrules"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/golang/mock/gomock"
)

func TestRuleTrigger(t *testing.T) {
	ctx := context.TODO()

	t.Run("Test invalid rule", func(t *testing.T) {
		// TODO
	})

	t.Run("execute rule, young age rule matched -> publish young rule execution requested", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, pubsub)

		// given
		ruleService.Put(ctx, youngAgeRule())

		// when
		defer func() {
			err := sut.Trigger(ctx, youngAgeRule())
			assert.NoError(t, err)
		}()

		// then
		pubsub.EXPECT().Publish(gomock.Any(), segmentrule.TriggerTopicName, segmentrule.RuleExecutionRequestedEvent{
			Rule: youngAgeRule(),
		}).Return(nil)
	})

	t.Run("execute rule, old age rule matched -> publish old rule execution requested", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, pubsub)

		// given
		ruleService.Put(ctx, oldAgeRule())

		// when
		defer func() {
			err := sut.Trigger(ctx, oldAgeRule())
			assert.NoError(t, err)
		}()

		// then
		pubsub.EXPECT().Publish(gomock.Any(), segmentrule.TriggerTopicName, segmentrule.RuleExecutionRequestedEvent{
			Rule: oldAgeRule(),
		}).Return(nil)
	})
}

func setup(t *testing.T) (*segmentrule.ManagementStub, *pubsub.MockPubsub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	ruleService := segmentrule.NewRuleManagementStub()
	pubsubMock := pubsub.NewMockPubsub(ctrl)

	return ruleService, pubsubMock, ctrl
}

func oldAgeRule() segmentrule.Spec {
	return predefinedrules.OldAgeEmailRule
}

func youngAgeRule() segmentrule.Spec {
	return predefinedrules.YoungAgeSmsRule
}

func nothingHappens() {}
