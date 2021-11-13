package usertaskexecutor

import (
	"context"
	"testing"

	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	"github.com/MarcGrol/userautomation/integrations/emailsending"
	"github.com/MarcGrol/userautomation/integrations/smssending"

	"github.com/MarcGrol/userautomation/coredata/predefinedrules"
	"github.com/MarcGrol/userautomation/coredata/predefinedusers"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestTaskExecution(t *testing.T) {
	ctx := context.TODO()

	t.Run("user-task execution requested, send email", func(t *testing.T) {
		// setup
		pubsub, smsSender, emailSender, reporter, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, reporter, smsSender, emailSender)

		// given

		// expect
		emailSender.EXPECT().Send(gomock.Any(), "marc@home.nl", "Your age is 50",
			"Hi Marc, your age is 50").Return(nil)

		// when
		err := sut.OnUserTaskExecutionRequestedEvent(ctx, usertask.UserTaskExecutionRequestedEvent{
			Task: usertask.Spec{
				RuleUID:    predefinedrules.OldAgeEmailRule.UID,
				ActionSpec: supportedactions.MailToOld,
				Reason:     usertask.ReasonUserAddedToSegment,
				User:       predefinedusers.Marc,
			},
		})
		assert.NoError(t, err)

		// then
		assert.Len(t, reporter.Reports, 1)
		assert.Equal(t, "Email with subject 'Your age is 50' has been sent to user 'marc@home.nl'", reporter.Reports[0].SuccessMessage)
	})

	t.Run("user-task execution requested, send sms", func(t *testing.T) {
		// setup
		pubsub, smsSender, emailSender, reporter, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub, reporter, smsSender, emailSender)

		// given

		// expect
		smsSender.EXPECT().Send(gomock.Any(), "+316333333", "Message to Pien").Return(nil)

		// when
		err := sut.OnUserTaskExecutionRequestedEvent(ctx, usertask.UserTaskExecutionRequestedEvent{
			Task: usertask.Spec{
				RuleUID:    predefinedrules.YoungAgeSmsRule.UID,
				ActionSpec: supportedactions.SmsToYoung,
				Reason:     usertask.ReasonUserAddedToSegment,
				User:       predefinedusers.Pien,
			},
		})
		assert.NoError(t, err)

		// then
		assert.Len(t, reporter.Reports, 1)
		assert.Equal(t, "Sms with content 'Message to Pien' has been sent to user '+316333333'", reporter.Reports[0].SuccessMessage)
	})
}

func setup(t *testing.T) (*pubsub.MockPubsub, *smssending.MockSmsSender, *emailsending.MockEmailSender, *usertask.ExecutionReporterStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	pubsubMock := pubsub.NewMockPubsub(ctrl)
	smsSender := smssending.NewMockSmsSender(ctrl)
	emailSender := emailsending.NewMockEmailSender(ctrl)
	reporter := usertask.NewExecutionReporterStub()

	return pubsubMock, smsSender, emailSender, reporter, ctrl
}
