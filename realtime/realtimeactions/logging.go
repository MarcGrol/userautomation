package realtimeactions

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

var logFunc = func(ctx context.Context, ruleName string, userStatus realtimecore.UserChangeStatus, oldState *realtimecore.User, newState *realtimecore.User) error {
	fmt.Printf("Action for rule '%s' triggered action om User '%s' - status: %+v\n",
		ruleName, getUserUid(oldState, newState), userStatus)
	return nil
}

func getUserUid(oldState *realtimecore.User, newState *realtimecore.User) string {
	if oldState != nil {
		return oldState.UID
	}
	return newState.UID
}

func LoggingAction() realtimecore.UserActionFunc {
	return logFunc
}
