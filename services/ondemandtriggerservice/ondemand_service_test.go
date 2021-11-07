package ondemandtriggerservice

import (
	"context"
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
		ruleService, userService, actionerMock, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, userService)

		// given
		noUsers()
		createOldAgeRule(ctx, t, ruleService, actionerMock)

		// when
		defer func() {
			err := sut.Trigger(ctx, "YoungRule")
			assert.Error(t, err)
		}()

		// then
		nothingHappens()
	})

	t.Run("execute rule, young age rule matched with no users", func(t *testing.T) {
		// setup
		ruleService, userService, actionerMock, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, userService)

		// given
		noUsers()
		createYoungAgeRule(ctx, t, ruleService, actionerMock)

		// when
		defer func() {
			err := sut.Trigger(ctx, "YoungRule")
			assert.NoError(t, err)
		}()

		// then
		nothingHappens()
	})

	t.Run("execute rule, young age rule matched -> trigger young rule execution", func(t *testing.T) {
		// setup
		ruleService, userService, actionerMock, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, userService)

		// given
		u := createUser(ctx, t, userService, 12)
		createYoungAgeRule(ctx, t, ruleService, actionerMock)

		// when
		defer sut.Trigger(ctx, "YoungRule")

		// then
		actionerMock.EXPECT().Perform(gomock.Any(), action.UserAction{
			RuleUID:  "YoungRule",
			Reason:   action.ReasonIsOnDemand,
			OldState: nil,
			NewState: &u,
		}).Return(nil)
	})

	t.Run("execute rule, old age rule matched -> no users", func(t *testing.T) {
		// setup
		ruleService, userService, actionerMock, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, userService)

		// given
		noUsers()
		createYoungAgeRule(ctx, t, ruleService, actionerMock)

		// when
		defer sut.Trigger(ctx, "OldRule")

		// then
		nothingHappens()
	})

	t.Run("execute rule, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, userService, actionerMock, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, userService)

		// given
		u := createUser(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService, actionerMock)

		// when
		defer sut.Trigger(ctx, "OldRule")

		// then
		actionerMock.EXPECT().Perform(gomock.Any(), action.UserAction{
			RuleUID:  "OldRule",
			Reason:   action.ReasonIsOnDemand,
			OldState: nil,
			NewState: &u,
		}).Return(nil)
	})

	t.Run("execute rule, old age rule matched multiple users -> 2 emails", func(t *testing.T) {
		// setup
		ruleService, userService, actionerMock, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(ruleService, userService)

		// given
		createUser(ctx, t, userService, 50)
		createOtherUser(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService, actionerMock)

		// when
		defer sut.Trigger(ctx, "OldRule")

		// then
		actionerMock.EXPECT().Perform(gomock.Any(), gomock.Any()).Return(nil).Times(2)
	})
}

func setup(t *testing.T) (*rule.SegmentRuleServiceStub, *user.UserServiceStub, *action.MockUserActioner, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	ruleService := rule.NewUserSegmentRuleServiceStub()
	userService := user.NewUserServiceStub()
	actionerMock := action.NewMockUserActioner(ctrl)

	return ruleService, userService, actionerMock, ctrl
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

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService rule.SegmentRuleService, actioner action.UserActioner) {
	err := segmentService.Put(ctx, rule.UserSegmentRule{
		UID: "OldRule",
		UserSegment: segment.UserSegment{
			UID:         "old users segment",
			Description: "old users segment",
			IsApplicableForUser: func(ctx context.Context, user user.User) (bool, error) {
				age, ok := user.Attributes["age"].(int)
				if !ok {
					return false, nil
				}
				return age > 40, nil
			},
		},
		Action: actioner,
	})
	if err != nil {
		t.Error(err)
	}
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService rule.SegmentRuleService, actioner action.UserActioner) {
	err := segmentService.Put(ctx, rule.UserSegmentRule{
		UID: "YoungRule",
		UserSegment: segment.UserSegment{
			UID:         "young users segment",
			Description: "young users segment",
			IsApplicableForUser: func(ctx context.Context, user user.User) (bool, error) {
				age, ok := user.Attributes["age"].(int)
				if !ok {
					return false, nil
				}
				return age < 18, nil
			},
		},
		Action:          actioner,
		AllowedTriggers: rule.TriggerUserChange,
	})
	if err != nil {
		t.Error(err)
	}
}
func nothingHappens() {}
