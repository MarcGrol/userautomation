package usertask

import (
	"context"
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
	UID        string
	RuleUID    string
	ActionSpec action.Spec
	Reason     Reason
	User       user.User
}

type UserTaskExecutionReport struct {
	TaskSpec       Spec
	Success        bool
	ErrorMessage   string
	SuccessMessage string
}

type ExecutionReporter interface {
	ReportExecution(ctx context.Context, report UserTaskExecutionReport) error
}
