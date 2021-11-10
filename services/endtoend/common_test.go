package endtoend

import (
	"context"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"testing"

	"github.com/MarcGrol/userautomation/core/rule"
	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/integrations/smssending"
	"github.com/golang/mock/gomock"
)

func setupSut(ctx context.Context) (rule.RuleService, user.Management) {
	sut := New(ctx)
	return sut.GetRuleService(), sut.GetUserService()
}

func createUser(ctx context.Context, t *testing.T, userService user.Management, age int) {
	err := userService.Put(ctx, user.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"first_name":    "Marc",
			"email_address": "marc@home.nl",
			"phone_number":  "+31611111111",
			"age":           age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func createOtherUser(ctx context.Context, t *testing.T, userService user.Management, age int) {
	err := userService.Put(ctx, user.User{
		UID: "2",
		Attributes: map[string]interface{}{
			"first_name":    "Eva",
			"email_address": "eva@home.nl",
			"phone_number":  "+31622222222",
			"age":           age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func modifyUser(ctx context.Context, t *testing.T, userService user.Management, age int) {
	err := userService.Put(ctx, user.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"first_name":    "Marc",
			"email_address": "marc@home.nl",
			"phone_number":  "+31611111111",
			"age":           age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func removeUser(ctx context.Context, t *testing.T, userService user.Management) {
	err := userService.Remove(ctx, "1")
	if err != nil {
		t.Error(err)
	}
}

func createOldAgeRule(ctx context.Context, t *testing.T, segmentService rule.RuleService,
	emailSender emailsending.EmailSender) {
	err := segmentService.Put(ctx, rule.RuleSpec{
		UID: "OldRule",
		SegmentSpec: segment.SegmentSpec{
			UID:            "old users segment",
			Description:    "old users segment",
			UserFilterName: user.FilterOldAge,
		},
		ActionSpec: supportedactions.MailToOld,
	})
	if err != nil {
		t.Error(err)
	}
}

func createYoungAgeRule(ctx context.Context, t *testing.T, segmentService rule.RuleService, smsSender smssending.SmsSender) {
	err := segmentService.Put(ctx, rule.RuleSpec{
		UID: "YoungRule",
		SegmentSpec: segment.SegmentSpec{
			UID:            "young users segment",
			Description:    "young users segment",
			UserFilterName: user.FilterYoungAge,
		},
		ActionSpec: supportedactions.SmsToYoung,
	})
	if err != nil {
		t.Error(err)
	}
}

func executeYoungAgeRuleReturnError(ctx context.Context, t *testing.T, ondemandService rule.TriggerRuleExecution) error {
	err := ondemandService.Trigger(ctx, "YoungRule")
	if err != nil {
		return err
	}
	return nil
}

func executeYoungAgeRule(ctx context.Context, t *testing.T, ondemandService rule.TriggerRuleExecution) {
	err := executeYoungAgeRuleReturnError(ctx, t, ondemandService)
	if err != nil {
		t.Error(err)
	}
}

func executeOldAgeRule(ctx context.Context, t *testing.T, ondemandService rule.TriggerRuleExecution) {
	err := ondemandService.Trigger(ctx, "OldRule")
	if err != nil {
		t.Error(err)
	}
}

func setupMocks(t *testing.T) (*emailsending.MockEmailSender, *smssending.MockSmsSender, *pubsub.MockPubsub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	mockEmailer := emailsending.NewMockEmailSender(ctrl)
	mockSmser := smssending.NewMockSmsSender(ctrl)
	ps := pubsub.NewMockPubsub(ctrl)
	return mockEmailer, mockSmser, ps, ctrl
}
