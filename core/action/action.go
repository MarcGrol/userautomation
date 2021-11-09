package action

import (
	"context"
	"github.com/MarcGrol/userautomation/core/usertask"
)

type ActionSpec struct {
	Name string
}

type ActionOnNamer interface {
	GetActionOnName(ctx context.Context, name string) (usertask.MockUserTaskExecutor, bool)
}