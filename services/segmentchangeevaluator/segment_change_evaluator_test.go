package segmentchangeevaluator

import (
	"context"
	"testing"

	"github.com/MarcGrol/userautomation/coredata/predefinedsegments"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/coredata/predefinedrules"
	"github.com/MarcGrol/userautomation/coredata/predefinedusers"
	"github.com/MarcGrol/userautomation/coredata/supportedattrs"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRuleEvaluation(t *testing.T) {
	ctx := context.TODO()

	t.Run("user added, no rules", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService)

		// given

		// when
		defer func() {
			err := sut.OnUserAddedToSegment(ctx, segment.UserAddedToSegmentEvent{
				SegmentUID: predefinedsegments.OldAgeSegment.UID,
				User:       predefinedusers.Marc,
			})
			assert.NoError(t, err)
		}()

		// then
		nothingHappens()
	})

	t.Run("user added, no rules match", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService)

		// given
		ruleService.Put(ctx, predefinedrules.YoungAgeSmsRule)

		// when
		defer func() {
			err := sut.OnUserAddedToSegment(ctx, segment.UserAddedToSegmentEvent{
				SegmentUID: predefinedsegments.OldAgeSegment.UID,
				User:       predefinedusers.Marc,
			})
			assert.NoError(t, err)
		}()

		// then
		nothingHappens()
	})

	t.Run("user added, 1 rule match", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService)

		// given
		ruleService.Put(ctx, predefinedrules.OldAgeEmailRule)
		ruleService.Put(ctx, predefinedrules.YoungAgeSmsRule)

		// when
		defer sut.OnUserAddedToSegment(ctx, segment.UserAddedToSegmentEvent{
			SegmentUID: predefinedsegments.OldAgeSegment.UID,
			User:       predefinedusers.Marc,
		})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
			Task: usertask.Spec{
				RuleUID:    predefinedrules.OldAgeEmailRule.UID,
				ActionSpec: predefinedrules.OldAgeEmailRule.ActionSpec,
				Reason:     usertask.ReasonUserAddedToSegment,
				User:       predefinedusers.Marc,
			},
		}).Return(nil)
	})

	t.Run("user added, 2 rules match", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService)

		// given
		ruleService.Put(ctx, predefinedrules.OldAgeEmailRule)
		ruleService.Put(ctx, predefinedrules.AotherOldAgeEmailRule)

		// when
		defer sut.OnUserAddedToSegment(ctx, segment.UserAddedToSegmentEvent{
			SegmentUID: predefinedsegments.OldAgeSegment.UID,
			User:       predefinedusers.Marc,
		})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), usertask.TopicName, gomock.Any()).Return(nil).Times(2)
	})

	t.Run("user-removed, nothing happens", func(t *testing.T) {
		// setup
		ruleService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService)

		// given
		ruleService.Put(ctx, predefinedrules.YoungAgeSmsRule)
		ruleService.Put(ctx, predefinedrules.OldAgeEmailRule)

		// when
		defer sut.OnUserRemovedFromSegment(ctx, segment.UserRemovedFromSegmentEvent{
			SegmentUID: predefinedrules.YoungAgeSmsRule.UID,
			User:       predefinedusers.Marc,
		})

		// then
		nothingHappens()
	})
}

func setup(t *testing.T) (*segmentrule.ManagementStub, *pubsub.MockPubsub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	ruleService := segmentrule.NewRuleManagementStub()
	pubsubMock := pubsub.NewMockPubsub(ctrl)

	return ruleService, pubsubMock, ctrl
}

func defaultUser() user.User {
	return predefinedusers.Marc
}

func createUserWithAge(ctx context.Context, t *testing.T, userService user.Management, age int) user.User {
	u := defaultUser()
	u.Attributes[supportedattrs.Age] = age
	err := userService.Put(ctx, u)
	if err != nil {
		t.Error(err)
	}
	return u
}

func otherUser() user.User {
	return predefinedusers.Eva
}

func createOtherUser(ctx context.Context, t *testing.T, userService user.Management, age int) {
	err := userService.Put(ctx, otherUser())
	if err != nil {
		t.Error(err)
	}
}

func oldAgeRule() segmentrule.Spec {
	return predefinedrules.OldAgeEmailRule
}

func createOldAgeRule(ctx context.Context, t *testing.T, ruleService segmentrule.Management) segmentrule.Spec {
	err := ruleService.Put(ctx, oldAgeRule())
	if err != nil {
		t.Error(err)
	}
	return oldAgeRule()
}

func youngAgeRule() segmentrule.Spec {
	return predefinedrules.YoungAgeSmsRule
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService segmentrule.Management) segmentrule.Spec {
	err := segmentService.Put(ctx, youngAgeRule())
	if err != nil {
		t.Error(err)
	}
	return youngAgeRule()
}
func nothingHappens() {}
