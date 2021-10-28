package rules

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/MarcGrol/userautomation/actions"
	"github.com/MarcGrol/userautomation/core"
	"github.com/MarcGrol/userautomation/userlookup"
)

func TestIt(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userLookup := userlookup.NewMockUserLookuper(ctrl)
	groupApi := actions.NewMockGroupApi(ctrl)

	// setup expectations
	userLookup.EXPECT().GetUserOnUid("123").Return(user, nil)
	groupApi.EXPECT().GroupExists("work.nl").Return(true, nil)
	groupApi.EXPECT().AddUserToGroup("work.nl", "123").Return(nil)

	sut := NewAddToGroupUserRule(userLookup, groupApi)

	err := core.EvaluateUserRule(sut, core.Event{
		EventName: "UserRegistered",
		UserUID:   "123",
		Payload:   map[string]interface{}{},
	})
	assert.NoError(t, err)
}

var user = core.User{
	UserUID:      "123",
	EmailAddress: "123@work.nl",
	PhoneNumber:  "+31612345678",
	CommunityUID: "xebia",
	Payload:      map[string]interface{}{},
}
