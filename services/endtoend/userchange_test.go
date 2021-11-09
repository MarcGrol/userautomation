package endtoend
//
//import (
//	"context"
//	"github.com/MarcGrol/userautomation/core/usertask"
//	"testing"
//
//	"github.com/MarcGrol/userautomation/core/rule"
//	"github.com/MarcGrol/userautomation/core/user"
//	"github.com/MarcGrol/userautomation/infra/pubsub"
//	"github.com/MarcGrol/userautomation/integrations/emailsending"
//	"github.com/MarcGrol/userautomation/integrations/smssending"
//
//	"github.com/golang/mock/gomock"
//)
//
//func TestUserChange(t *testing.T) {
//	testCases := []struct {
//		name  string
//		given func(c givenContext)
//		when  func(c whenContext)
//		then  func(tc thenContext)
//	}{
//		{
//			name:  "create user, no rule exists",
//			given: nothingGiven(),
//			when: func(c whenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//			},
//			then: nothingHappens(),
//		},
//		{
//			name: "create user, no rule matched",
//			given: func(c givenContext) {
//				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
//			},
//			when: func(c whenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//			},
//			then: nothingHappens(),
//		},
//		{
//			name: "create user, young age rule matched -> sms",
//			given: func(c givenContext) {
//				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
//			},
//			when: func(c whenContext) {
//				createUser(c.ctx, t, c.userService, 12)
//			},
//			then: func(c thenContext) {
//				c.ps.EXPECT().Publish(gomock.Any(), usertask.TopicName,
//					usertask.UserTaskExecutionRequestedEvent{Task: usertask.UserTask{
//						RuleUID:  "",
//						Reason:   0,
//						OldState: nil,
//						NewState: nil,
//					},
//					})
//			},
//		},
//
//		{
//			name: "create user, old age rule matched -> email",
//			given: func(c givenContext) {
//				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
//			},
//			when: func(c whenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//			},
//			then: func(c thenContext) {
//				c.ps.EXPECT().Publish(gomock.Any(), usertask.TopicName,
//					usertask.UserTaskExecutionRequestedEvent{Task: usertask.UserTask{
//						RuleUID:  "",
//						Reason:   0,
//						OldState: nil,
//						NewState: nil,
//					},
//					})
//			},
//		},
//		{
//			name: "modify user, no rule exist",
//			given: func(c givenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//			},
//			when: func(c whenContext) {
//				modifyUser(c.ctx, t, c.userService, 12)
//			},
//			then: nothingHappens(),
//		},
//		{
//			name: "modify user, no rule matched",
//			given: func(c givenContext) {
//				createUser(c.ctx, t, c.userService, 12)
//				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
//			},
//			when: func(c whenContext) {
//				modifyUser(c.ctx, t, c.userService, 14)
//			},
//			then: nothingHappens(),
//		},
//		{
//			name: "modify user, young age rule matched -> sms",
//			given: func(c givenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
//			},
//			when: func(c whenContext) {
//				modifyUser(c.ctx, t, c.userService, 12)
//			},
//			then: func(c thenContext) {
//				c.ps.EXPECT().Publish(gomock.Any(), usertask.TopicName,
//					usertask.UserTaskExecutionRequestedEvent{Task: usertask.UserTask{
//						RuleUID:  "",
//						Reason:   0,
//						OldState: nil,
//						NewState: nil,
//					},
//					})
//			},
//		},
//		{
//			name: "modify user, old age rule matched -> email",
//			given: func(c givenContext) {
//				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
//				createUser(c.ctx, t, c.userService, 12)
//			},
//			when: func(c whenContext) {
//				modifyUser(c.ctx, t, c.userService, 50)
//			},
//			then: func(c thenContext) {
//				c.ps.EXPECT().Publish(gomock.Any(), usertask.TopicName,
//					usertask.UserTaskExecutionRequestedEvent{Task: usertask.UserTask{
//						RuleUID:  "",
//						Reason:   0,
//						OldState: nil,
//						NewState: nil,
//					},
//					})
//			},
//		},
//		{
//			name: "modify user, remains young",
//			given: func(c givenContext) {
//				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
//				createUser(c.ctx, t, c.userService, 12)
//
//			},
//			when: func(c whenContext) {
//				modifyUser(c.ctx, t, c.userService, 14)
//			},
//			then: nothingHappens(),
//		},
//		{
//			name:  "delete user, no user exists",
//			given: nothingGiven(),
//			when: func(c whenContext) {
//				removeUser(c.ctx, t, c.userService)
//			},
//			then: nothingHappens(),
//		},
//
//		{
//			name: "delete user, no rule exist",
//			given: func(c givenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//			},
//			when: func(c whenContext) {
//				removeUser(c.ctx, t, c.userService)
//			},
//			then: nothingHappens(),
//		},
//		{
//			name: "delete user, no rule matched",
//			given: func(c givenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//				createYoungAgeRule(c.ctx, t, c.ruleService, c.smser)
//
//			},
//			when: func(c whenContext) {
//				defer removeUser(c.ctx, t, c.userService)
//
//			},
//			then: nothingHappens(),
//		},
//		{
//			name: "delete user, young age rule matched",
//			given: func(c givenContext) {
//				createUser(c.ctx, t, c.userService, 50)
//				createOldAgeRule(c.ctx, t, c.ruleService, c.emailer)
//			},
//			when: func(c whenContext) {
//				defer removeUser(c.ctx, t, c.userService)
//
//			},
//			then: nothingHappens(),
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			// Start with a fully fresh environment each test
//			ctx := context.TODO()
//			mockEmailer, mockSmser, mockPubsub, ctrl := setupMocks(t)
//			defer ctrl.Finish()
//
//			ruleService, userService := setupSut(ctx)
//
//			// execute the test
//			tc.given(givenContext{
//				ctx:         ctx,
//				ruleService: ruleService,
//				userService: userService,
//				emailer:     mockEmailer,
//				smser:       mockSmser,
//			})
//
//			defer tc.when(whenContext{
//				ctx:         ctx,
//				userService: userService,
//			})
//
//			tc.then(thenContext{
//				emailer: mockEmailer,
//				smser:   mockSmser,
//				ps: mockPubsub,
//			})
//		})
//	}
//}
//
//type givenContext struct {
//	ctx         context.Context
//	ruleService rule.RuleService
//	userService user.Management
//	emailer     emailsending.EmailSender
//	smser       smssending.SmsSender
//}
//
//type whenContext struct {
//	ctx         context.Context
//	ruleService rule.RuleService
//	userService user.Management
//}
//
//type thenContext struct {
//	emailer *emailsending.MockEmailSender
//	smser   *smssending.MockSmsSender
//	ps *pubsub.MockPubsub
//}
//
//func nothingGiven() func(c givenContext) {
//	return func(c givenContext) {}
//}
//
//func nothingHappens() func(c thenContext) {
//	return func(c thenContext) {}
//}
