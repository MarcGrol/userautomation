package segmentrule

import (
	"context"
)

type ManagementStub struct {
	Rules map[string]Spec
}

func NewRuleManagementStub() *ManagementStub {
	return &ManagementStub{
		Rules: map[string]Spec{},
	}
}

func (s *ManagementStub) Put(ctx context.Context, r Spec) error {
	s.Rules[r.UID] = r
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
