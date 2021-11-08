package endtoend

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

func setupSut(ctx context.Context) (rule.SegmentRuleService, user.Service) {
	sut := New(ctx)
	return sut.GetRuleService(), sut.GetUserService()
}

func createUser(ctx context.Context, t *testing.T, userService user.Service, age int) {
	err := userService.Put(ctx, user.User{
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

func modifyUser(ctx context.Context, t *testing.T, userService user.Service, age int) {
	err := userService.Put(ctx, user.User{
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

func removeUser(ctx context.Context, t *testing.T, userService user.Service) {
	err := userService.Remove(ctx, "1")
	if err != nil {
		t.Error(err)
	}
}

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService rule.SegmentRuleService,
	emailSender emailsending.EmailSender) {
	err := segmentService.Put(ctx, rule.UserSegmentRule{
		UID: "OldRule",
		UserSegment: segment.UserSegment{
			UID:            "old users segment",
			Description:    "old users segment",
			UserFilterName: user.FilterOldAge,
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
			UserFilterName: user.FilterYoungAge,
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

func executeYoungAgeRule(ctx context.Context, t *testing.T, ondemandService rule.SegmentRuleExecutionTrigger) {
	err := executeYoungAgeRuleReturnError(ctx, t, ondemandService)
	if err != nil {
		t.Error(err)
	}
}

func executeOldAgeRule(ctx context.Context, t *testing.T, ondemandService rule.SegmentRuleExecutionTrigger) {
	err := ondemandService.Trigger(ctx, "OldRule")
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
