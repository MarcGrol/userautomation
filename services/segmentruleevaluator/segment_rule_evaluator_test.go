package segmentruleevaluator

import (
	"context"
	"testing"

	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	. "github.com/MarcGrol/userautomation/coredata/supportedattrs"
	"github.com/MarcGrol/userautomation/coredata/supportedrules"
	"github.com/MarcGrol/userautomation/infra/pubsub"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/segmentrule"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRuleEvaluation(t *testing.T) {
	ctx := context.TODO()

	t.Run("execute non-existing rule", func(t *testing.T) {
		// setup
		ruleService, userService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService, userService)

		// given
		noUsers()
		createOldAgeRule(ctx, t, ruleService)

		// when
		defer func() {
			err := sut.OnRuleExecutionRequestedEvent(ctx, segmentrule.RuleExecutionRequestedEvent{Rule: youngAgeRule()})
			assert.Error(t, err)
		}()

		// then
		nothingHappens()
	})

	t.Run("execute rule, young age rule matched with no users", func(t *testing.T) {
		// setup
		ruleService, userService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService, userService)

		// given
		noUsers()
		createYoungAgeRule(ctx, t, ruleService)

		// when
		defer func() {
			err := sut.OnRuleExecutionRequestedEvent(ctx, segmentrule.RuleExecutionRequestedEvent{Rule: youngAgeRule()})
			assert.NoError(t, err)
		}()

		// then
		nothingHappens()
	})

	t.Run("execute rule, young age rule matched -> trigger young rule execution", func(t *testing.T) {
		// setup
		ruleService, userService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService, userService)

		// given
		u := createUser(ctx, t, userService, 12)
		r := createYoungAgeRule(ctx, t, ruleService)

		// when
		defer sut.OnRuleExecutionRequestedEvent(ctx, segmentrule.RuleExecutionRequestedEvent{Rule: youngAgeRule()})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
			Task: usertask.Spec{
				ActionSpec: r.ActionSpec,
				Reason:     usertask.ReasonSegmentRuleExecuted,
				User:       u,
			},
		}).Return(nil)
	})

	t.Run("execute rule, old age rule matched -> no users", func(t *testing.T) {
		// setup
		ruleService, userService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService, userService)

		// given
		noUsers()
		createYoungAgeRule(ctx, t, ruleService)

		// when
		defer sut.OnRuleExecutionRequestedEvent(ctx, segmentrule.RuleExecutionRequestedEvent{Rule: oldAgeRule})

		// then
		nothingHappens()
	})

	t.Run("execute rule, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, userService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService, userService)

		// given
		u := createUser(ctx, t, userService, 50)
		r := createOldAgeRule(ctx, t, ruleService)

		// when
		defer sut.OnRuleExecutionRequestedEvent(ctx, segmentrule.RuleExecutionRequestedEvent{Rule: oldAgeRule})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
			Task: usertask.Spec{
				ActionSpec: r.ActionSpec,
				Reason:     usertask.ReasonSegmentRuleExecuted,
				User:       u,
			},
		}).Return(nil)
	})

	t.Run("execute rule, old age rule matched multiple users -> 2 emails", func(t *testing.T) {
		// setup
		ruleService, userService, pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, ruleService, userService)

		// given
		createUser(ctx, t, userService, 50)
		createOtherUser(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService)

		// when
		defer sut.OnRuleExecutionRequestedEvent(ctx, segmentrule.RuleExecutionRequestedEvent{Rule: oldAgeRule})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), usertask.TopicName, gomock.Any()).Return(nil).Times(2)
	})
}

func setup(t *testing.T) (*segmentrule.ManagementStub, *user.UserManagementStub, *pubsub.MockPubsub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	ruleService := segmentrule.NewRuleManagementStub()
	userService := user.NewUserManagementStub()
	pubsubMock := pubsub.NewMockPubsub(ctrl)

	return ruleService, userService, pubsubMock, ctrl
}

func noUsers() {}

func createUser(ctx context.Context, t *testing.T, userService user.Management, age int) user.User {
	u := user.User{
		UID: "1",
		Attributes: map[string]interface{}{
			FirstName:    "Marc",
			EmailAddress: "marc@home.nl",
			PhoneNumber:  "+31611111111",
			Age:          age,
		},
	}
	err := userService.Put(ctx, u)
	if err != nil {
		t.Error(err)
	}
	return u
}

func createOtherUser(ctx context.Context, t *testing.T, userService user.Management, age int) {
	err := userService.Put(ctx, user.User{
		UID: "2",
		Attributes: map[string]interface{}{
			FirstName:    "Eva",
			EmailAddress: "eva@home.nl",
			PhoneNumber:  "+31622222222",
			Age:          age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

var oldAgeRule = segmentrule.Spec{
	UID: "OldRule",
	SegmentSpec: segment.Spec{
		UID:            "old users segment",
		Description:    "old users segment",
		UserFilterName: user.FilterOldAgeName,
	},
	ActionSpec: supportedactions.MailToOld,
}

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService segmentrule.Service) segmentrule.Spec {
	err := segmentService.Put(ctx, oldAgeRule)
	if err != nil {
		t.Error(err)
	}
	return oldAgeRule
}

func youngAgeRule() segmentrule.Spec {
	return supportedrules.YoungAgeSmsRule
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService segmentrule.Service) segmentrule.Spec {
	err := segmentService.Put(ctx, youngAgeRule())
	if err != nil {
		t.Error(err)
	}
	return youngAgeRule()
}
func nothingHappens() {}
