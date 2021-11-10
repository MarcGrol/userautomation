package ondemandtriggerservice

import (
	"context"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	"testing"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOnDemand(t *testing.T) {
	ctx := context.TODO()

	t.Run("execute non-existing rule", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, pubsub)

		// given
		createRule(ctx, t, ruleService, oldAgeRule)

		// when
		defer func() {
			err := sut.Trigger(ctx, "YoungRule")
			assert.Error(t, err)
		}()

		// then
		nothingHappens()
	})

	t.Run("execute rule, young age rule matched -> publish young rule execution requested", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, pubsub)

		// given
		createRule(ctx, t, ruleService, youngAgeRule)

		// when
		defer sut.Trigger(ctx, "YoungRule")

		// then
		pubsub.EXPECT().Publish(gomock.Any(), segmentrule.TriggerTopicName, segmentrule.RuleExecutionRequestedEvent{
			Rule: youngAgeRule,
		}).Return(nil)
	})

	t.Run("execute rule, old age rule matched -> publish old rule execution requested", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, pubsub)

		// given
		createRule(ctx, t, ruleService, oldAgeRule)

		// when
		defer sut.Trigger(ctx, "OldRule")

		// then
		pubsub.EXPECT().Publish(gomock.Any(), segmentrule.TriggerTopicName, segmentrule.RuleExecutionRequestedEvent{
			Rule: oldAgeRule,
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

var oldAgeRule = segmentrule.Spec{
	UID: "OldRule",
	SegmentSpec: segment.Spec{
		UID:            "old users segment",
		Description:    "old users segment",
		UserFilterName: user.FilterOldAge,
	},
	ActionSpec: supportedactions.MailToOld,
}

var youngAgeRule = segmentrule.Spec{
	UID: "YoungRule",
	SegmentSpec: segment.Spec{
		UID:            "young users segment",
		Description:    "young users segment",
		UserFilterName: user.FilterYoungAge,
	},
	ActionSpec: supportedactions.SmsToYoung,
}

func nothingHappens() {}
