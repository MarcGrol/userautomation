package rules

import (
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

	// setup expectations
	userLookup.EXPECT().GetUserOnQuery("publishCount > 10 && loginCount > 20").Return([]core.User{
		{
			UserUID:      "123",
			EmailAddress: "123@work.nl",
			PhoneNumber:  "+31612345678",
			CommunityUID: "xebia",
			Payload: map[string]interface{}{
				"FirstName": "Marc",
			},
		},
	}, nil)
	emailer.EXPECT().Send("123@work.nl", "We praise your activity", "Hi Marc, well done").Return(nil)

	sut := NewPraiseActiveUserRule(userLookup, emailer)

	err := core.EvaluateUserRule(sut, core.Event{
		EventName: "Timer",
		Payload:   map[string]interface{}{},
	})
	assert.NoError(t, err)
}
