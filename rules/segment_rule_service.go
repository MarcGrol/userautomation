package rules

import (
	"context"
	"sync"
)

type userSegmentRuleService struct {
	sync.Mutex
	rules map[string]UserSegmentRule
}

func NewUserSegmentRuleService() SegmentRuleService {
	return &userSegmentRuleService{
		rules: map[string]UserSegmentRule{},
	}
}

func (s *userSegmentRuleService) Put(ctx context.Context, SegmentRule UserSegmentRule) error {
	s.Lock()
	defer s.Unlock()

	s.rules[SegmentRule.Name] = SegmentRule

	return nil
}

func (s *userSegmentRuleService) Get(ctx context.Context, segmentUID string) (UserSegmentRule, bool, error) {
	s.Lock()
	defer s.Unlock()

	SegmentRule, exists := s.rules[segmentUID]

	return SegmentRule, exists, nil
}

func (s *userSegmentRuleService) Delete(ctx context.Context, segmentUID string) error {
	s.Lock()
	defer s.Unlock()

	delete(s.rules, segmentUID)

	return nil
}

func (s *userSegmentRuleService) List(ctx context.Context) ([]UserSegmentRule, error) {
	s.Lock()
	defer s.Unlock()

	segments := []UserSegmentRule{}

	for _, sd := range s.rules {
		segments = append(segments, sd)
	}

	return segments, nil
}
