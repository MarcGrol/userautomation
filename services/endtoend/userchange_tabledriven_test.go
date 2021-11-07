package endtoend

import (
	"context"
	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/user"
	emailsending2 "github.com/MarcGrol/userautomation/integrations/emailsending"
	smssending2 "github.com/MarcGrol/userautomation/integrations/smssending"
	"testing"

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
			then: nothingHappens(),
		},
		{
			name: "create user, no rule matched",
			given: func(c givenContext) {
				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
			},
			when: func(c whenContext) {
				createUser(c.ctx, t, c.userService, 50)
			},
			then: nothingHappens(),
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
			then: nothingHappens(),
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
			then: nothingHappens(),
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
			then: nothingHappens(),
		},
		{
			name:  "delete user, no user exists",
			given: nothingGiven(),
			when: func(c whenContext) {
				removeUser(c.ctx, t, c.userService)
			},
			then: nothingHappens(),
		},

		{
			name: "delete user, no rule exist",
			given: func(c givenContext) {
				createUser(c.ctx, t, c.userService, 50)
			},
			when: func(c whenContext) {
				removeUser(c.ctx, t, c.userService)
			},
			then: nothingHappens(),
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
			then: nothingHappens(),
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
			then: nothingHappens(),
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

type givenContext struct {
	ctx         context.Context
	ruleService rule.SegmentRuleService
	userService user.Service
	emailer     emailsending2.EmailSender
	smser       smssending2.SmsSender
}

type whenContext struct {
	ctx         context.Context
	ruleService rule.SegmentRuleService
	userService user.Service
}

type thenContext struct {
	emailer *emailsending2.MockEmailSender
	smser   *smssending2.MockSmsSender
}

func nothingGiven() func(c givenContext) {
	return func(c givenContext) {}
}

func nothingHappens() func(c thenContext) {
	return func(c thenContext) {}
}
