package segmentruletrigger

import (
	"context"
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
		createRule(ctx, t, ruleService, youngAgeRule())

		// when
		defer sut.Trigger(ctx, youngAgeRule())

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
		createRule(ctx, t, ruleService, oldAgeRule())

		// when
		defer sut.Trigger(ctx, oldAgeRule())

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

func createRule(ctx context.Context, t *testing.T, segmentService segmentrule.Service, r segmentrule.Spec) {
	err := segmentService.Put(ctx, r)
	if err != nil {
		t.Error(err)
	}
}

func oldAgeRule() segmentrule.Spec {
	return predefinedrules.OldAgeEmailRule
}

func youngAgeRule() segmentrule.Spec {
	return predefinedrules.YoungAgeSmsRule
}

func nothingHappens() {}
