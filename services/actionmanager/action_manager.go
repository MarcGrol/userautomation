package actionmanager

import (
	"context"
	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
)

type service struct {
}

func New() action.ActionManager {
	return &service{}
}

var actionMap = map[string]action.Spec{
	supportedactions.SmsToYoungName: supportedactions.SmsToYoung,
	supportedactions.MailToOldName:  supportedactions.MailToOld,
}

func (m *service) GetActionSpecOnName(ctx context.Context, name string) (action.Spec, bool, error) {
	a, exists := actionMap[name]

	return a, exists, nil
}

func (m *service) ListActionSpecs(ctx context.Context) ([]action.Spec, error) {
	actions := []action.Spec{}
	for _, a := range actionMap {
		actions = append(actions, a)
	}
	return actions, nil
}
