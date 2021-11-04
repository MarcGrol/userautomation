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

func TestAddToGroup(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userLookup := userlookup.NewMockUserLookuper(ctrl)
	groupApi := batchactions.NewMockGroupApi(ctrl)

	testCases := []struct {
		name              string
		event             batchcore.Event
		setupExpectations func()
		expectedResult    error
	}{
		{
			name: "Unsupported event",
			event: batchcore.Event{
				EventName: "Timer",
				Payload:   map[string]interface{}{},
			},
			setupExpectations: func() {},
			expectedResult:    nil,
		},
		{
			name: "Supported event",
			event: batchcore.Event{
				EventName: "UserRegistered",
				UserUID:   "123",
				Payload:   map[string]interface{}{},
			},
			setupExpectations: func() {
				userLookup.EXPECT().GetUserOnUid(gomock.Any(), "123").Return(testUser, nil)
				groupApi.EXPECT().GroupExists(gomock.Any(), "work.nl").Return(true, nil)
				groupApi.EXPECT().AddUserToGroup(gomock.Any(), "work.nl", "123").Return(nil)
			},
			expectedResult: nil,
		},
	}
	for _, tc := range testCases {
		sut := NewAddToGroupUserRule(userLookup, groupApi)
		t.Run(tc.name, func(t *testing.T) {
			tc.setupExpectations()
			err := batchcore.EvaluateUserRule(context.TODO(), sut, tc.event)
			assert.NoError(t, err)
		})
	}
}