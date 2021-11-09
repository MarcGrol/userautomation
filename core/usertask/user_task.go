package usertask

import (
	"context"
	"fmt"

	"github.com/MarcGrol/userautomation/core/user"
)

//go:generate mockgen -source=user_task.go -destination=task_executor_mock.go -package=usertask UserTaskExecutor
type UserTaskExecutor interface {
	Perform(ctx context.Context, task UserTask) error
}

type ReasonForAction int

const (
	ReasonIsUserAddedToSegment ReasonForAction = iota
	ReasonIsOnDemand
	ReasonIsCron
)

type UserTask struct {
	RuleUID  string
	Reason   ReasonForAction
	NewState user.User
}

func (t UserTask) String() string {
	return fmt.Sprintf("UserTaskExecutor for rule '%s' triggered action om User '%s' - status: %+v\n",
		t.RuleUID, t.NewState, t.Reason)
}