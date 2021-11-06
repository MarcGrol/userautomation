package allwiredtogether

import (
	"context"
	"github.com/MarcGrol/userautomation/segments"
	"testing"

	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/integrations/smssending"
	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/useractions/emailaction"
	"github.com/MarcGrol/userautomation/useractions/sms"
	"github.com/MarcGrol/userautomation/users"
	"github.com/golang/mock/gomock"
)

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
	emailSender emailsending.EmailSender) {
	err := segmentService.Put(ctx, rules.UserSegmentRule{
		Name: "OldRule",
		UserSegment: segments.UserSegment{
			Name: "old users segment",
			IsApplicableForUser: func(ctx context.Context, user users.User) (bool, error) {
				age, ok := user.Attributes["age"].(int)
				if !ok {
					return false, nil
				}
				return age > 40, nil
			},
			Users: []users.User{}, // start empty
		},
		PerformActionForUser: emailaction.EmailerAction("old rule fired", "Hoi {{.firstname}}, your age is {{.age}}", emailSender),
	})
	if err != nil {
		t.Error(err)
	}
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService rules.SegmentRuleService, smsSender smssending.SmsSender) {
	err := segmentService.Put(ctx, rules.UserSegmentRule{
		Name: "YoungRule",
		UserSegment: segments.UserSegment{
			Name: "young users segment",
			IsApplicableForUser: func(ctx context.Context, user users.User) (bool, error) {
				age, ok := user.Attributes["age"].(int)
				if !ok {
					return false, nil
				}
				return age < 18, nil
			},
			Users: []users.User{}, // start empty
		},
		PerformActionForUser: sms.SmsAction("young rule fired for {{.firstname}}: your age is {{.age}}", smsSender),
	})
	if err != nil {
		t.Error(err)
	}
}

func setupMocks(t *testing.T) (*emailsending.MockEmailSender, *smssending.MockSmsSender, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockEmailer := emailsending.NewMockEmailSender(ctrl)
	mockSmser := smssending.NewMockSmsSender(ctrl)
	return mockEmailer, mockSmser, ctrl
}