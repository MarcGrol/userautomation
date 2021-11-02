package rules

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/core"
	"github.com/MarcGrol/userautomation/userlookup"
)

func TestPraiseActiveUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userLookup := userlookup.NewMockUserLookuper(ctrl)
	emailer := actions.NewMockEmailer(ctrl)

	testCases := []struct {
		name              string
		event             core.Event
		setupExpectations func()
		expectedResult    error
	}{
		{
			name: "Unsupported event",
			event: core.Event{
				EventName: "UserRegistered",
				Payload:   map[string]interface{}{},
			},
			setupExpectations: func() {},
			expectedResult:    nil,
		},
		{
			name: "Supported event",
			event: core.Event{
				EventName: "Timer",
				Payload:   map[string]interface{}{},
			},
			setupExpectations: func() {
				userLookup.EXPECT().GetUserOnQuery(gomock.Any(), "publishCount > 10 && loginCount > 20").Return([]core.User{testUser}, nil)
				emailer.EXPECT().Send(gomock.Any(), "123@work.nl", "We praise your activity", "Hi Marc, well done").Return(nil)
			},
			expectedResult: nil,
		},
	}
	for _, tc := range testCases {
		sut := NewPraiseActiveUserRule(userLookup, emailer)
		t.Run(tc.name, func(t *testing.T) {
			tc.setupExpectations()
			err := core.EvaluateUserRule(context.TODO(), sut, tc.event)
			assert.NoError(t, err)
		})
	}
}
