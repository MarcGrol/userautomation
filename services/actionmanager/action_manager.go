package actionmanager

import (
	"context"

	"github.com/MarcGrol/userautomation/core/action"
	"github.com/MarcGrol/userautomation/coredata/supportedactions"
	"github.com/MarcGrol/userautomation/infra/datastore"
)

type service struct {
	store datastore.Datastore
}

func New(store datastore.Datastore) action.Management {
	store.EnforceDataType(action.Spec{})
	return &service{
		store: store,
	}
}

func (m *service) GetActionSpecOnName(ctx context.Context, name string) (action.Spec, bool, error) {
	found, exists, err := m.store.Get(ctx, name)
	if err != nil {
		return action.Spec{}, false, err
	}
	if !exists {
		return action.Spec{}, false, nil
	}
	return found.(action.Spec), true, nil
}

func (m *service) ListActionSpecs(ctx context.Context) ([]action.Spec, error) {
	items, err := m.store.GetAll(ctx)
	if err != nil {
		return []action.Spec{}, err
	}
	actions := []action.Spec{}
	for _, i := range items {
		actions = append(actions, i.(action.Spec))
	}
	return actions, nil
}

func (m *service) Preprov(ctx context.Context) error {
	err := m.store.Put(ctx, supportedactions.SmsToYoungName, supportedactions.SmsToYoung)
	if err != nil {
		return err
	}

	err = m.store.Put(ctx, supportedactions.MailToOldName, supportedactions.MailToOld)
	if err != nil {
		return err
	}
	return nil
}
