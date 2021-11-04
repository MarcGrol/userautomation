package realtimeactions

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

var logFunc = func(ctx context.Context, ruleName string, userStatus realtimecore.UserStatus, user realtimecore.User) error {
	fmt.Printf("Action for rule '%s' triggered action om User '%s' - status: %+v\n", ruleName, user.UID, userStatus)
	return nil
}

func LoggingAction() realtimecore.UserActionFunc {
	return logFunc
}
