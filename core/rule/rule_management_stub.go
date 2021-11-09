package rule

import (
	"context"
)

type RuleManagementStub struct {
	Rules map[string]RuleSpec
}

func NewRuleManagementStub() *RuleManagementStub {
	return &RuleManagementStub{
		Rules: map[string]RuleSpec{},
	}
}

func (s *RuleManagementStub) Put(ctx context.Context, r RuleSpec) error {
	s.Rules[r.UID] = r
	return nil
}

func (s *RuleManagementStub) Remove(ctx context.Context, uid string) error {
	delete(s.Rules, uid)
	return nil
}

func (s *RuleManagementStub) Get(ctx context.Context, uid string) (RuleSpec, bool, error) {
	item, exists := s.Rules[uid]
	return item, exists, nil
}

func (s *RuleManagementStub) List(ctx context.Context) ([]RuleSpec, error) {
	items := []RuleSpec{}
	for _, i := range s.Rules {
		items = append(items, i)
	}
	return items, nil
}
