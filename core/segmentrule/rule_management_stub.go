package segmentrule

import (
	"context"

	"github.com/MarcGrol/userautomation/core/util"
)

type ManagementStub struct {
	util.NoPreProvNeeded
	util.NoWebNeeded
	Rules map[string]Spec
}

func NewRuleManagementStub() *ManagementStub {
	return &ManagementStub{
		Rules: map[string]Spec{},
	}
}

func (s *ManagementStub) Put(ctx context.Context, rule Spec) error {
	s.Rules[rule.UID] = rule
	return nil
}

func (s *ManagementStub) Remove(ctx context.Context, uid string) error {
	delete(s.Rules, uid)
	return nil
}

func (s *ManagementStub) Get(ctx context.Context, uid string) (Spec, bool, error) {
	item, exists := s.Rules[uid]
	return item, exists, nil
}

func (s *ManagementStub) List(ctx context.Context) ([]Spec, error) {
	items := []Spec{}
	for _, i := range s.Rules {
		items = append(items, i)
	}
	return items, nil
}
