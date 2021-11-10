package usertask

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/segmentrule"

	"github.com/MarcGrol/userautomation/core/user"
)

//go:generate mockgen -source=user_task.go -destination=task_executor_mock.go -package=usertask UserTaskExecutor
type UserTaskExecutor interface {
	Perform(ctx context.Context, task Spec) error
}

type ReasonForAction int

const (
	ReasonUserAddedToSegment ReasonForAction = iota
	ReasonSegmentRuleExecuted
	ReasonUserRuleExecuted
)

type Spec struct {
	RuleSpec segmentrule.Spec
	Reason   ReasonForAction
	User     user.User
}

func (t Spec) String() string {
	return fmt.Sprintf("UserTaskExecutor for rule '%s' triggered action om User '%s' - status: %+v\n",
		t.RuleSpec.UID, t.User, t.Reason)
}
