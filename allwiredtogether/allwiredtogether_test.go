package allwiredtogether

import (
	"context"
	"testing"

	"github.com/MarcGrol/userautomation/actions/email"
	"github.com/MarcGrol/userautomation/actions/sms"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/users"
	"github.com/golang/mock/gomock"
)

func TestUsingTableStrategy(t *testing.T) {
	testCases := []struct {
		name  string
		given func(c givenContext)
		when  func(c whenContext)
		then  func(tc thenContext)
	}{
		{
			name:  "create user, no rule exists",
			given: nothingGiven(),
			when: func(c whenContext) {
				createUser(c.ctx, t, c.userService, 50)
			},
			then: nothing(),
		},
		{
			name: "create user, no rule matched",
			given: func(c givenContext) {
				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
			},
			when: func(c whenContext) {
				createUser(c.ctx, t, c.userService, 50)
			},
			then: nothing(),
		},
		{
			name: "create user, young age rule matched -> sms",
			given: func(c givenContext) {
				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
			},
			when: func(c whenContext) {
				createUser(c.ctx, t, c.userService, 12)
			},
			then: func(c thenContext) {
				c.smser.EXPECT().Send(gomock.Any(), "+31611111111",
					"young rule fired for Marc: your age is 12").Return(nil)
			},
		},

		{
			name: "create user, old age rule matched -> email",
			given: func(c givenContext) {
				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
			},
			when: func(c whenContext) {
				createUser(c.ctx, t, c.userService, 50)
			},
			then: func(c thenContext) {
				c.emailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
					"old rule fired", "Hoi Marc, your age is 50").Return(nil)
			},
		},
		{
			name: "modify user, no rule exist",
			given: func(c givenContext) {
				createUser(c.ctx, t, c.userService, 50)
			},
			when: func(c whenContext) {
				modifyUser(c.ctx, t, c.userService, 12)
			},
			then: nothing(),
		},
		{
			name: "modify user, no rule matched",
			given: func(c givenContext) {
				createUser(c.ctx, t, c.userService, 12)
				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
			},
			when: func(c whenContext) {
				modifyUser(c.ctx, t, c.userService, 14)
			},
			then: nothing(),
		},
		{
			name: "modify user, young age rule matched -> sms",
			given: func(c givenContext) {
				createUser(c.ctx, t, c.userService, 50)
				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
			},
			when: func(c whenContext) {
				modifyUser(c.ctx, t, c.userService, 12)
			},
			then: func(c thenContext) {
				c.smser.EXPECT().Send(gomock.Any(), "+31611111111",
					"young rule fired for Marc: your age is 12").Return(nil)
			},
		},
		{
			name: "modify user, old age rule matched -> email",
			given: func(c givenContext) {
				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
				createUser(c.ctx, t, c.userService, 12)
			},
			when: func(c whenContext) {
				modifyUser(c.ctx, t, c.userService, 50)
			},
			then: func(c thenContext) {
				c.emailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
					"old rule fired", "Hoi Marc, your age is 50").Return(nil)
			},
		},
		{
			name: "modify user, remains young",
			given: func(c givenContext) {
				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
				createUser(c.ctx, t, c.userService, 12)

			},
			when: func(c whenContext) {
				modifyUser(c.ctx, t, c.userService, 14)
			},
			then: nothing(),
		},
		{
			name:  "delete user, no user exists",
			given: nothingGiven(),
			when: func(c whenContext) {
				removeUser(c.ctx, t, c.userService)
			},
			then: nothing(),
		},

		{
			name: "delete user, no rule exist",
			given: func(c givenContext) {
				createUser(c.ctx, t, c.userService, 50)
			},
			when: func(c whenContext) {
				removeUser(c.ctx, t, c.userService)
			},
			then: nothing(),
		},
		{
			name: "delete user, no rule matched",
			given: func(c givenContext) {
				createUser(c.ctx, t, c.userService, 50)
				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)

			},
			when: func(c whenContext) {
				defer removeUser(c.ctx, t, c.userService)

			},
			then: nothing(),
		},
		{
			name: "delete user, young age rule matched",
			given: func(c givenContext) {
				createUser(c.ctx, t, c.userService, 50)
				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
			},
			when: func(c whenContext) {
				defer removeUser(c.ctx, t, c.userService)

			},
			then: nothing(),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Start with a fully fresh environment each test
			ctx := context.TODO()
			mockEmailer, mockSmser, ctrl := setupMocks(t)
			defer ctrl.Finish()

			ruleService, userService := setupSut(ctx)

			// execute the test
			tc.given(givenContext{
				ctx:         ctx,
				ruleService: ruleService,
				userService: userService,
				emailer:     mockEmailer,
				smser:       mockSmser,
			})

			defer tc.when(whenContext{
				ctx:         ctx,
				userService: userService,
			})

			tc.then(thenContext{
				emailer: mockEmailer,
				smser:   mockSmser,
			})
		})
	}
}

func TestUsingSubTests(t *testing.T) {
	ctx := context.TODO()

	t.Run("create user, no rule exists", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given

		// when
		defer createUser(ctx, t, userService, 50)

		// then
	})

	t.Run("create user, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer createUser(ctx, t, userService, 50)

		// then
	})

	t.Run("create user, young age rule matched -> sms", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer createUser(ctx, t, userService, 12)

		// then
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111",
			"young rule fired for Marc: your age is 12").Return(nil)
	})

	t.Run("create user, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer createUser(ctx, t, userService, 50)

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
			"old rule fired", "Hoi Marc, your age is 50").Return(nil)

	})

	t.Run("modify user, no rule exist", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// expect

		// when
		defer modifyUser(ctx, t, userService, 12)

		// then
		createUser(ctx, t, userService, 50)

	})

	t.Run("modify user, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 12)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer modifyUser(ctx, t, userService, 14)

		// then

	})

	t.Run("modify user, young age rule matched -> sms", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer modifyUser(ctx, t, userService, 12)

		// then
		mockSmser.EXPECT().Send(gomock.Any(), "+31611111111",
			"young rule fired for Marc: your age is 12").Return(nil)

	})

	t.Run("modify user, old age rule matched -> email", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createUser(ctx, t, userService, 12)

		// when
		defer modifyUser(ctx, t, userService, 50)

		// then
		mockEmailer.EXPECT().Send(gomock.Any(), "marc@home.nl",
			"old rule fired", "Hoi Marc, your age is 50").Return(nil)

	})

	t.Run("modify user, remains young", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createOldAgeRule(ctx, t, ruleService, mockEmailer)
		createUser(ctx, t, userService, 12)

		// when
		defer modifyUser(ctx, t, userService, 14)

		// then
	})

	t.Run("delete user, no user exists", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given

		// when
		defer removeUser(ctx, t, userService)

		// then
	})

	t.Run("delete user, no rule exist", func(t *testing.T) {
		// setup
		_, userService := setupSut(ctx)
		_, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)

		// when
		defer removeUser(ctx, t, userService)

		// then
	})

	t.Run("delete user, no rule matched", func(t *testing.T) {
		// setup
		ruleService, userService  := setupSut(ctx)
		_, mockSmser, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)
		createYoungAgeRule(ctx, t, ruleService, mockSmser)

		// when
		defer removeUser(ctx, t, userService)

		// then
	})

	t.Run("delete user, young age rule matched", func(t *testing.T) {
		// setup
		ruleService, userService := setupSut(ctx)
		mockEmailer, _, ctrl := setupMocks(t)
		defer ctrl.Finish()

		// given
		createUser(ctx, t, userService, 50)
		createOldAgeRule(ctx, t, ruleService, mockEmailer)

		// when
		defer removeUser(ctx, t, userService)

		// then
	})
}

func setupSut(ctx context.Context) (rules.SegmentRuleService, users.UserService) {
	sut := New(ctx)
	return sut.GetSegmentRuleService(), sut.GetUserService()
}

func createUser(ctx context.Context, t *testing.T, userService users.UserService, age int) {
	err := userService.Put(ctx, users.User{
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

func modifyUser(ctx context.Context, t *testing.T, userService users.UserService, age int) {
	err := userService.Put(ctx, users.User{
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

func removeUser(ctx context.Context, t *testing.T, userService users.UserService) {
	err := userService.Remove(ctx, "1")
	if err != nil {
		t.Error(err)
	}
}

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService rules.SegmentRuleService,
	emailSender email.EmailSender) {
	err := segmentService.Put(ctx, rules.UserSegmentRule{
		Name: "OldRule",
		IsApplicableForUser: func(ctx context.Context, user users.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age > 40, nil
		},
		PerformAction: email.EmailerAction("old rule fired", "Hoi {{.firstname}}, your age is {{.age}}", emailSender),
	})
	if err != nil {
		t.Error(err)
	}
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService rules.SegmentRuleService, smsSender sms.SmsSender) {
	err := segmentService.Put(ctx, rules.UserSegmentRule{
		Name: "YoungRule",
		IsApplicableForUser: func(ctx context.Context, user users.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age < 18, nil
		},
		PerformAction: sms.SmsAction("young rule fired for {{.firstname}}: your age is {{.age}}", smsSender),
	})
	if err != nil {
		t.Error(err)
	}
}

func setupMocks(t *testing.T) (*email.MockEmailSender, *sms.MockSmsSender, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockEmailer := email.NewMockEmailSender(ctrl)
	mockSmser := sms.NewMockSmsSender(ctrl)
	return mockEmailer, mockSmser, ctrl
}

type givenContext struct {
	ctx         context.Context
	ruleService rules.SegmentRuleService
	userService users.UserService
	emailer     email.EmailSender
	smser   sms.SmsSender
}

type whenContext struct {
	ctx         context.Context
	ruleService rules.SegmentRuleService
	userService users.UserService
}

type thenContext struct {
	emailer *email.MockEmailSender
	smser   *sms.MockSmsSender
}

func nothingGiven() func(c givenContext) {
	return func(c givenContext) {}
}

func nothing() func(c thenContext) {
	return func(c thenContext) {}
}
