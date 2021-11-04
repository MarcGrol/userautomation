package batchrules

import (
	"context"
	"github.com/MarcGrol/userautomation/batch/batchactions"
	"github.com/MarcGrol/userautomation/batch/batchcore"
	"github.com/MarcGrol/userautomation/batch/userlookup"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPraiseActiveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userLookup := userlookup.NewMockUserLookuper(ctrl)
	emailer := batchactions.NewMockEmailer(ctrl)

	testCases := []struct {
		name              string
		event             batchcore.Event
		setupExpectations func()
		expectedResult    error
	}{
		{
			name: "Unsupported event",
			event: batchcore.Event{
				EventName: "UserRegistered",
				Payload:   map[string]interface{}{},
			},
			setupExpectations: func() {},
			expectedResult:    nil,
		},
		{
			name: "Supported event",
			event: batchcore.Event{
				EventName: "Timer",
				Payload:   map[string]interface{}{},
			},
			setupExpectations: func() {
				userLookup.EXPECT().GetUserOnQuery(gomock.Any(), "publishCount > 10 && loginCount > 20").Return([]batchcore.User{testUser}, nil)
				emailer.EXPECT().Send(gomock.Any(), "123@work.nl", "We praise your activity", "Hi Marc, well done").Return(nil)
			},
			expectedResult: nil,
		},
	}
	for _, tc := range testCases {
		sut := NewPraiseActiveUserRule(userLookup, emailer)
		t.Run(tc.name, func(t *testing.T) {
			tc.setupExpectations()
			err := batchcore.EvaluateUserRule(context.TODO(), sut, tc.event)
			assert.NoError(t, err)
		})
	}
}
