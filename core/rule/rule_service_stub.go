package rule

import (
	"context"
)

type SegmentRuleServiceStub struct {
	Rules map[string]UserSegmentRule
}

func NewUserSegmentRuleServiceStub() *SegmentRuleServiceStub {
	return &SegmentRuleServiceStub{
		Rules: map[string]UserSegmentRule{},
	}
}

func (s *SegmentRuleServiceStub) Put(ctx context.Context, r UserSegmentRule) error {
	s.Rules[r.UID] = r
	return nil
}

func (s *SegmentRuleServiceStub) Remove(ctx context.Context, uid string) error {
	delete(s.Rules, uid)
	return nil
}

func (s *SegmentRuleServiceStub) Get(ctx context.Context, uid string) (UserSegmentRule, bool, error) {
	item, exists := s.Rules[uid]
	return item, exists, nil
}

func (s *SegmentRuleServiceStub) List(ctx context.Context) ([]UserSegmentRule, error) {
	items := []UserSegmentRule{}
	for _, i := range s.Rules {
		items = append(items, i)
	}
	return items, nil
}
