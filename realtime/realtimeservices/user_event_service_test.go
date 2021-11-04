package realtimeservices

import (
	"context"
	"github.com/golang/mock/gomock"
	"testing"

	"github.com/MarcGrol/userautomation/realtime/realtimeactions"
	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

func TestCreateUser(t *testing.T) {
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
		createMarc(ctx, t, userService)

	})

	t.Run("create user, no rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect

		// when
		createMarc(ctx, t, userService)

	})

	t.Run("create user, young age rule matched -> sms", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect
		mockSmser.EXPECT().Send(gomock.Any(), "+31633333333", "young rule fired for Freek: your age is 12").Return(nil)

		// when
		createFreek(ctx, t, userService)
	})

	t.Run("create user, old age rule matched -> email", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// expect
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl", "old rule fired", "Hoi Marc, your age is 50").Return(nil)

		// when
		createMarc(ctx, t, userService)

	})

	t.Run("modify user, no user exists", func(t *testing.T) {
		_, userService := setup(ctx)

		// given

		// expect

		// when
		adjustMarc(ctx, t, userService)

	})

	t.Run("modify user, no rule exist", func(t *testing.T) {
		_, userService := setup(ctx)

		// given
		createMarc(ctx, t, userService)

		// expect

		// when
		adjustMarc(ctx, t, userService)

	})

	t.Run("modify user, no rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createFreek(ctx, t, userService)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// expect

		// when
		adjustFreek(ctx, t, userService)

	})

	t.Run("modify user, young age rule matched -> sms", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createMarc(ctx, t, userService)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111", "young rule fired for Marc: your age is 10").Return(nil)

		// when
		adjustMarc(ctx, t, userService)

	})

	t.Run("modify user, old age rule matched -> email", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createFreek(ctx, t, userService)

		// expect
		mockEmailer.EXPECT().Send(gomock.Any(), "freek@home.nl", "old rule fired", "Hoi Freek, your age is 41").Return(nil)

		// when
		adjustFreekAgain(ctx, t, userService)

	})

	t.Run("modify user, remains young", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createFreek(ctx, t, userService)

		// expect

		// when
		adjustFreek(ctx, t, userService)

	})

	t.Run("delete user, no user exists", func(t *testing.T) {
		_, userService := setup(ctx)

		// given

		// expect

		// when
		deleteMarc(ctx, t, userService)

	})

	t.Run("delete user, no rule exist", func(t *testing.T) {
		_, userService := setup(ctx)

		// given
		createMarc(ctx, t, userService)

		// expect

		// when
		deleteMarc(ctx, t, userService)

	})

	t.Run("delete user, no rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createMarc(ctx, t, userService)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// expect

		// when
		deleteMarc(ctx, t, userService)

	})

	t.Run("delete user, young age rule matched", func(t *testing.T) {
		ruleService, userService := setup(ctx)

		// given
		createMarc(ctx, t, userService)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// expect

		// when
		deleteMarc(ctx, t, userService)

	})
}

func createMarc(ctx context.Context, t *testing.T, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          50, // old
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func adjustMarc(ctx context.Context, t *testing.T, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          10, // now young
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func deleteMarc(ctx context.Context, t *testing.T, userService realtimecore.UserService) {
	err := userService.Remove(ctx, "1")
	if err != nil {
		t.Error(err)
	}
}

func createEva(ctx context.Context, t *testing.T, userService realtimecore.UserService) {

	err := userService.Put(ctx, realtimecore.User{
		UID: "2",
		Attributes: map[string]interface{}{
			"firstname":    "Eva",
			"emailaddress": "eva@home.nl",
			"phonenumber":  "+31622222222",
			"age":          48, // old
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func createFreek(ctx context.Context, t *testing.T, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID: "3",
		Attributes: map[string]interface{}{
			"firstname":    "Freek",
			"emailaddress": "freek@home.nl",
			"phonenumber":  "+31633333333",
			"age":          12, // young
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func adjustFreek(ctx context.Context, t *testing.T, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID: "3",
		Attributes: map[string]interface{}{
			"firstname":    "Freek",
			"emailaddress": "freek@home.nl",
			"phonenumber":  "+31633333333",
			"age":          13, // slightly older but still young
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func adjustFreekAgain(ctx context.Context, t *testing.T, userService realtimecore.UserService) {
	err := userService.Put(ctx, realtimecore.User{
		UID: "3",
		Attributes: map[string]interface{}{
			"firstname":    "Freek",
			"emailaddress": "freek@home.nl",
			"phonenumber":  "+31633333333",
			"age":          41, // big increase in age, now old
		},
	})
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
