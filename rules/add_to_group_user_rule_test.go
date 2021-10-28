package rules

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"github.com/golang/mock/gomock"

	"github.com/MarcGrol/userautomation/userlookup"

)

func TestIt(  t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	userLookup := userlookup.NewMockUserLookuper(ctrl)
	groupApi := actions.NewGroupApi(ctrl)

	// setup expectations
	userLookup.EXPECT().GetUserOnUid().Return()

	sut := NewAddToGroupUserRule(userLookup, groupApi),

	event := api.Event{
		EventName: "Timer",
		Payload:   map[string]interface{}{},
	}

	err := EvaluateUserRule(sut, event)
	assert.NoError(t, err)

}
