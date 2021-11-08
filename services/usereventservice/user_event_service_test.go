package usereventservice

import (
	"context"
	"testing"

	"github.com/MarcGrol/userautomation/actions/emailaction"
	"github.com/MarcGrol/userautomation/actions/smsaction"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/integrations/smssending"
	"github.com/golang/mock/gomock"
)

func TestUsingClassicSubTests(t *testing.T) {
	ctx := context.TODO()

	t.Run("user-created user, no rule exists", func(t *testing.T) {
		// setup
		_, sut := setup()
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		noUsers()
		noRules()

		// when
		defer userCreated(ctx, t, sut, 50)

		// then
		nothingHappens()
	})

	t.Run("used-created user, no rule matched", func(t *testing.T) {
		// setup
		ruleService, sut := setup()
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		noUsers()
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer userCreated(ctx, t, sut, 50)

		// then
		nothingHappens()

	})

	t.Run("user-created user, young age rule matched -> sms", func(t *testing.T) {
		// setup
		ruleService, userEventService := setup()
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		noUsers()
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer userCreated(ctx, t, userEventService, 12)

		// then
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111",
			"young rule fired for Marc: your age is 12").Return(nil)
	})

	t.Run("user-created, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, sut := setup()
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		noUsers()
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer userCreated(ctx, t, sut, 50)

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
			"old rule fired", "Hoi Marc, your age is 50").Return(nil)

	})

	t.Run("user-modified, no rule exist", func(t *testing.T) {
		// setup
		_, userService := setup()
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		userCreated(ctx, t, userService, 50)
		noRules()

		// when
		defer userModified(ctx, t, userService, 50, 12)

		// then
		nothingHappens()

	})

	t.Run("user-modified, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setup()
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		userCreated(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer userModified(ctx, t, userService, 50, 51)

		// then
		nothingHappens()

	})

	t.Run("user-modified, young age rule matched -> sms", func(t *testing.T) {
		// setup
		ruleService, userService := setup()
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		userCreated(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer userModified(ctx, t, userService, 50, 12)

		// then
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111",
			"young rule fired for Marc: your age is 12").Return(nil)

	})

	t.Run("user-modified, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, userService := setup()
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		userCreated(ctx, t, userService, 12)

		// when
		defer userModified(ctx, t, userService, 12, 50)

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
			"old rule fired", "Hoi Marc, your age is 50").Return(nil)

	})

	t.Run("user-modified, remains young", func(t *testing.T) {
		// setup
		ruleService, userService := setup()
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		userCreated(ctx, t, userService, 12)

		// when
		defer userModified(ctx, t, userService, 12, 14)

		// then
		nothingHappens()

	})

	t.Run("user-removed, no user exists", func(t *testing.T) {
		// setup
		_, userService := setup()
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		noUsers()
		noRules()

		// when
		defer userRemoved(ctx, t, userService)

		// then
		nothingHappens()

	})

	t.Run("user-removed, no rule exist", func(t *testing.T) {
		// setup
		_, userService := setup()
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		userCreated(ctx, t, userService, 50)
		noRules()

		// when
		defer userRemoved(ctx, t, userService)

		// then
		nothingHappens()

	})

	t.Run("user-removed, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setup()
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		userCreated(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer userRemoved(ctx, t, userService)

		// then
		nothingHappens()

	})

	t.Run("user-removed, young age rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setup()
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		userCreated(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer userRemoved(ctx, t, userService)

		// then
		nothingHappens()

	})
}

func setup() (rule.SegmentRuleService, UserEventService) {
	ruleService := rule.NewUserSegmentRuleServiceStub()
	return ruleService, NewUserEventService(nil, ruleService)
}

func noUsers() {}

func userCreated(ctx context.Context, t *testing.T, service UserEventService, age int) {
	err := service.OnUserCreated(ctx, user.CreatedEvent{
		UserState: user.User{
			UID: "1",
			Attributes: map[string]interface{}{
				"firstname":    "Marc",
				"emailaddress": "marc@home.nl",
				"phonenumber":  "+31611111111",
				"age":          age,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
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

func userModified(ctx context.Context, t *testing.T, service UserEventService, oldAge, nnewAge int) {
	err := service.OnUserModified(ctx, user.ModifiedEvent{
		OldUserState: user.User{
			UID: "1",
			Attributes: map[string]interface{}{
				"firstname":    "Marc",
				"emailaddress": "marc@home.nl",
				"phonenumber":  "+31611111111",
				"age":          oldAge,
			},
		},
		NewUserState: user.User{
			UID: "1",
			Attributes: map[string]interface{}{
				"firstname":    "Marc",
				"emailaddress": "marc@home.nl",
				"phonenumber":  "+31611111111",
				"age":          nnewAge,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func userRemoved(ctx context.Context, t *testing.T, service UserEventService) {
	err := service.OnUserRemoved(ctx, user.RemovedEvent{
		UserState: user.User{
			UID: "1",
			Attributes: map[string]interface{}{
				"firstname":    "Marc",
				"emailaddress": "marc@home.nl",
				"phonenumber":  "+31611111111",
				"age":          50,
			},
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func noRules() {}

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService rule.SegmentRuleService,
	emailSender emailsending.EmailSender) {
	err := segmentService.Put(ctx, rule.UserSegmentRule{
		UID: "OldRule",
		UserSegment: segment.UserSegment{
			UID:            "old users segment",
			Description:    "old users segment",
			UserFilterName: segment.FilterOldAge,
		},
		Action: emailaction.NewEmailAction("old rule fired", "Hoi {{.firstname}}, your age is {{.age}}", emailSender),
	})
	if err != nil {
		t.Error(err)
	}
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService rule.SegmentRuleService, smsSender smssending.SmsSender) {
	err := segmentService.Put(ctx, rule.UserSegmentRule{
		UID: "YoungRule",
		UserSegment: segment.UserSegment{
			UID:            "young users segment",
			Description:    "young users segment",
			UserFilterName: segment.FilterYoungAge,
		},
		Action:          smsaction.New("young rule fired for {{.firstname}}: your age is {{.age}}", smsSender),
		AllowedTriggers: rule.TriggerUserChange,
	})
	if err != nil {
		t.Error(err)
	}
}

func executeYoungAgeRuleReturnError(ctx context.Context, t *testing.T, ondemandService rule.SegmentRuleExecutionTrigger) error {
	err := ondemandService.Trigger(ctx, "YoungRule")
	if err != nil {
		return err
	}
	return nil
}

func setupMocks(t *testing.T) (*emailsending.MockEmailSender, *smssending.MockSmsSender, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockEmailer := emailsending.NewMockEmailSender(ctrl)
	mockSmser := smssending.NewMockSmsSender(ctrl)
	return mockEmailer, mockSmser, ctrl
}

func nothingHappens() {}
