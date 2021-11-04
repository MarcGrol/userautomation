package realtimeservices

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/MarcGrol/userautomation/realtime/realtimeactions"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

func TestIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockEmailer := realtimeactions.NewMockEmailer(ctrl)
	mockSmser := realtimeactions.NewMockSmsSender(ctrl)

	ctx := context.TODO()

	t.Run("create user, no rule exists", func(t *testing.T) {
		_, userService := setup(ctx)

		// given

		// expect

		// when
		createUser(ctx, t, userService, 50)

	})

	t.Run("create user, no rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect

		// when
		createUser(ctx, t, userService, 50)

	})

	t.Run("create user, young age rule matched -> sms", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111", "young rule fired for Marc: your age is 12").Return(nil)

		// when
		createUser(ctx, t, userService, 12)
	})

	t.Run("create user, old age rule matched -> email", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// expect
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl", "old rule fired", "Hoi Marc, your age is 50").Return(nil)

		// when
		createUser(ctx, t, userService, 50)

	})

	t.Run("modify user, no rule exist", func(t *testing.T) {
		_, userService := setup(ctx)

		// given
		createUser(ctx, t, userService, 50)

		// expect

		// when
		modifyUser(ctx, t, userService, 12)

	})

	t.Run("modify user, no rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createUser(ctx, t, userService, 12)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// expect

		// when
		modifyUser(ctx, t, userService, 14)

	})

	t.Run("modify user, young age rule matched -> sms", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createUser(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111", "young rule fired for Marc: your age is 12").Return(nil)

		// when
		modifyUser(ctx, t, userService, 12)

	})

	t.Run("modify user, old age rule matched -> email", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createUser(ctx, t, userService, 12)

		// expect
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl", "old rule fired", "Hoi Marc, your age is 50").Return(nil)

		// when
		modifyUser(ctx, t, userService, 50)

	})

	t.Run("modify user, remains young", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createUser(ctx, t, userService, 12)

		// expect

		// when
		modifyUser(ctx, t, userService, 14)

	})

	t.Run("delete user, no user exists", func(t *testing.T) {
		_, userService := setup(ctx)

		// given

		// expect

		// when
		removeUser(ctx, t, userService)

	})

	t.Run("delete user, no rule exist", func(t *testing.T) {
		_, userService := setup(ctx)

		// given
		createUser(ctx, t, userService, 50)

		// expect

		// when
		removeUser(ctx, t, userService)

	})

	t.Run("delete user, no rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createUser(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect

		// when
		removeUser(ctx, t, userService)

	})

	t.Run("delete user, young age rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createUser(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// expect

		// when
		removeUser(ctx, t, userService)

	})
}

func createUser(ctx context.Context, t *testing.T, userService realtimecore.UserService, age int) {
	err := userService.Put(ctx, realtimecore.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func modifyUser(ctx context.Context, t *testing.T, userService realtimecore.UserService, age int) {
	err := userService.Put(ctx, realtimecore.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func removeUser(ctx context.Context, t *testing.T, userService realtimecore.UserService) {
	err := userService.Remove(ctx, "1")
	if err != nil {
		t.Error(err)
	}
}

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService realtimecore.SegmentRuleService,
	emailSender realtimeactions.Emailer) {
	err := segmentService.Put(ctx, realtimecore.UserSegmentRule{
		Name: "OldRule",
		IsApplicableForUser: func(ctx context.Context, user realtimecore.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age > 40, nil
		},
		PerformAction: realtimeactions.EmailerAction("old rule fired", "Hoi {{.firstname}}, your age is {{.age}}", emailSender),
	})
	if err != nil {
		t.Error(err)
	}
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService realtimecore.SegmentRuleService, smsSender realtimeactions.SmsSender) {
	err := segmentService.Put(ctx, realtimecore.UserSegmentRule{
		Name: "YoungRule",
		IsApplicableForUser: func(ctx context.Context, user realtimecore.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age < 18, nil
		},
		PerformAction: realtimeactions.SmsAction("young rule fired for {{.firstname}}: your age is {{.age}}", smsSender),
	})
	if err != nil {
		t.Error(err)
	}
}

func setup(ctx context.Context) (realtimecore.SegmentRuleService, realtimecore.UserService) {
	pubsub := NewPubSub()

	ruleService := NewUserSegmentRuleService()

	userEventService := NewUserEventService(pubsub, ruleService)
	userEventService.Subscribe(ctx)

	userService := NewUserService(pubsub)

	return ruleService, userService
}
