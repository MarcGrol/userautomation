package ondemandexecutor

import (
	"context"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"testing"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestOnDemand(t *testing.T) {
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
			err := sut.OnRuleExecutionRequestedEvent(ctx, rule.RuleExecutionRequestedEvent{Rule: youngAgeRule})
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
			err := sut.OnRuleExecutionRequestedEvent(ctx, rule.RuleExecutionRequestedEvent{Rule: youngAgeRule})
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
		createYoungAgeRule(ctx, t, ruleService)

		// when
		defer sut.OnRuleExecutionRequestedEvent(ctx, rule.RuleExecutionRequestedEvent{Rule: youngAgeRule})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), action.TopicName, action.ActionExecutionRequestedEvent{
			Action: action.UserAction{
				RuleUID:  "YoungRule",
				Reason:   action.ReasonIsOnDemand,
				OldState: nil,
				NewState: &u,
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
		defer sut.OnRuleExecutionRequestedEvent(ctx, rule.RuleExecutionRequestedEvent{Rule: oldAgeRule})

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
		createOldAgeRule(ctx, t, ruleService)

		// when
		defer sut.OnRuleExecutionRequestedEvent(ctx, rule.RuleExecutionRequestedEvent{Rule: oldAgeRule})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), action.TopicName, action.ActionExecutionRequestedEvent{
			Action: action.UserAction{
				RuleUID:  "OldRule",
				Reason:   action.ReasonIsOnDemand,
				OldState: nil,
				NewState: &u,
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
		defer sut.OnRuleExecutionRequestedEvent(ctx, rule.RuleExecutionRequestedEvent{Rule: oldAgeRule})

		// then
		pubsub.EXPECT().Publish(gomock.Any(), action.TopicName, gomock.Any()).Return(nil).Times(2)
	})
}

func setup(t *testing.T) (*rule.SegmentRuleServiceStub, *user.UserServiceStub, *pubsub.MockPubsub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	ruleService := rule.NewUserSegmentRuleServiceStub()
	userService := user.NewUserServiceStub()
	pubsubMock := pubsub.NewMockPubsub(ctrl)

	return ruleService, userService, pubsubMock, ctrl
}

func noUsers() {}

func createUser(ctx context.Context, t *testing.T, userService user.Service, age int) user.User {
	u := user.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          age,
		},
	}
	err := userService.Put(ctx, u)
	if err != nil {
		t.Error(err)
	}
	return u
}

func createOtherUser(ctx context.Context, t *testing.T, userService user.Service, age int) {
	err := userService.Put(ctx, user.User{
		UID: "2",
		Attributes: map[string]interface{}{
			"firstname":    "Eva",
			"emailaddress": "eva@home.nl",
			"phonenumber":  "+31622222222",
			"age":          age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

var oldAgeRule = rule.UserSegmentRule{
	UID: "OldRule",
	UserSegment: segment.UserSegment{
		UID:            "old users segment",
		Description:    "old users segment",
		UserFilterName: user.FilterOldAge,
	},
	Action: nil,
}

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService rule.SegmentRuleService) {
	err := segmentService.Put(ctx, oldAgeRule)
	if err != nil {
		t.Error(err)
	}
}

var youngAgeRule = rule.UserSegmentRule{
	UID: "YoungRule",
	UserSegment: segment.UserSegment{
		UID:            "young users segment",
		Description:    "young users segment",
		UserFilterName: user.FilterYoungAge,
	},
	Action:          nil,
	AllowedTriggers: rule.TriggerUserChange,
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService rule.SegmentRuleService) {
	err := segmentService.Put(ctx, youngAgeRule)
	if err != nil {
		t.Error(err)
	}
}
func nothingHappens() {}
