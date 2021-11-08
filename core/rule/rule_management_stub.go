package rule

import (
	"context"
)

type SegmentRuleManagementStub struct {
	Rules map[string]UserSegmentRule
}

func NewUserSegmentRuleManagementStub() *SegmentRuleManagementStub {
	return &SegmentRuleManagementStub{
		Rules: map[string]UserSegmentRule{},
	}
}

func (s *SegmentRuleManagementStub) Put(ctx context.Context, r UserSegmentRule) error {
	s.Rules[r.UID] = r
	return nil
}

func (s *SegmentRuleManagementStub) Remove(ctx context.Context, uid string) error {
	delete(s.Rules, uid)
	return nil
}

func (s *SegmentRuleManagementStub) Get(ctx context.Context, uid string) (UserSegmentRule, bool, error) {
	item, exists := s.Rules[uid]
	return item, exists, nil
}

func (s *SegmentRuleManagementStub) List(ctx context.Context) ([]UserSegmentRule, error) {
	items := []UserSegmentRule{}
	for _, i := range s.Rules {
		items = append(items, i)
	}
	return items, nil
}
