package usertask

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/core/user"
)

//go:generate mockgen -source=user_task.go -destination=task_executor_mock.go -package=usertask UserTaskExecutor
type UserTaskExecutor interface {
	Perform(ctx context.Context, task Spec) (string, error)
}

type Reason int

const (
	ReasonUserAddedToSegment Reason = iota
	ReasonSegmentRuleExecuted
	ReasonUserRuleExecuted
)

type Spec struct {
	RuleUID    string
	ActionSpec action.Spec
	Reason     Reason
	User       user.User
}

func (t Spec) String() string {
	return fmt.Sprintf("UserTaskExecutor triggered action %s on User '%s' - status: %+v\n",
		t.ActionSpec.Name, t.User, t.Reason)
}

type UserTaskExecutionReport struct {
	TaskSpec     Spec
	Success      bool
	ErrorMessage string
	SuccessMessage string
}

type ExecutionReporter interface {
	ReportExecution(ctx context.Context, report UserTaskExecutionReport) error
}
