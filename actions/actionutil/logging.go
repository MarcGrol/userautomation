package actionutil

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/rules"
	"github.com/MarcGrol/userautomation/users"
)

var LogFunc = func(ctx context.Context, ruleName string, userStatus rules.UserChangeStatus, oldState *users.User, newState *users.User) error {
	fmt.Printf("Action for rule '%s' triggered action om User '%s' - status: %+v\n",
		ruleName, getUserUid(oldState, newState), userStatus)
	return nil
}

func getUserUid(oldState *users.User, newState *users.User) string {
	if oldState != nil {
		return oldState.UID
	}
	return newState.UID
}

func LoggingAction() rules.UserActionFunc {
	return LogFunc
}
