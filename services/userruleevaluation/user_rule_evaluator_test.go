package userruleevaluator

import (
	"context"
	"testing"

	"github.com/MarcGrol/userautomation/core/userrule"
	"github.com/MarcGrol/userautomation/core/usertask"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"

	"github.com/MarcGrol/userautomation/coredata/predefinedrules"
	"github.com/MarcGrol/userautomation/coredata/predefinedusers"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestRuleEvaluation(t *testing.T) {
	ctx := context.TODO()

	t.Run("user-rule execution requested, fire task", func(t *testing.T) {
		// setup
		pubsub, ctrl := setup(t)
		defer ctrl.Finish()
		sut := New(pubsub)

		// given

		// when
		defer func() {
			err := sut.OnRuleExecutionRequestedEvent(ctx, userrule.RuleExecutionRequestedEvent{
				Rule: userrule.Spec{
					UID:         "1",
					Description: "My test",
					User:        predefinedusers.Marc,
					ActionSpec:  supportedactions.MailToOld,
				},
			})
			assert.NoError(t, err)
		}()

		// then
		pubsub.EXPECT().Publish(gomock.Any(), usertask.TopicName, usertask.UserTaskExecutionRequestedEvent{
			Task: usertask.Spec{
				RuleUID:    "",
				ActionSpec: predefinedrules.OldAgeEmailRule.ActionSpec,
				Reason:     usertask.ReasonUserRuleExecuted,
				User:       predefinedusers.Marc,
			},
		}).Return(nil)
	})
}

func setup(t *testing.T) (*pubsub.MockPubsub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	pubsubMock := pubsub.NewMockPubsub(ctrl)

	return pubsubMock, ctrl
}
