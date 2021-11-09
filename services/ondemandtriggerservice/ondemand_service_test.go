package ondemandtriggerservice

import (
	"context"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"testing"

	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
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
		pubsub.EXPECT().Publish(gomock.Any(), rule.TriggerTopicName, rule.RuleExecutionRequestedEvent{
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
		pubsub.EXPECT().Publish(gomock.Any(), rule.TriggerTopicName, rule.RuleExecutionRequestedEvent{
			Rule: oldAgeRule,
		}).Return(nil)
	})
}

func setup(t *testing.T) (*rule.RuleManagementStub, *pubsub.MockPubsub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	ruleService := rule.NewRuleManagementStub()
	pubsubMock := pubsub.NewMockPubsub(ctrl)

	return ruleService, pubsubMock, ctrl
}

func createRule(ctx context.Context, t *testing.T, segmentService rule.RuleService, r rule.RuleSpec) {
	err := segmentService.Put(ctx, r)
	if err != nil {
		t.Error(err)
	}
}

var oldAgeRule = rule.RuleSpec{
	UID: "OldRule",
	SegmentSpec: segment.SegmentSpec{
		UID:            "old users segment",
		Description:    "old users segment",
		UserFilterName: user.FilterOldAge,
	},
	ActionSpec: action.ActionSpec{Name: "email"},
}

var youngAgeRule = rule.RuleSpec{
	UID: "YoungRule",
	SegmentSpec: segment.SegmentSpec{
		UID:            "young users segment",
		Description:    "young users segment",
		UserFilterName: user.FilterYoungAge,
	},
	ActionSpec: action.ActionSpec{Name: "sms"},
}

func nothingHappens() {}
