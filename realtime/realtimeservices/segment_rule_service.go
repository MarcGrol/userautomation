package realtimeservices

import (
	"context"
	"sync"

	"github.com/MarcGrol/userautomation/realtime/realtimecore"
)

type userSegmentRuleService struct {
	sync.Mutex
	rules map[string]realtimecore.UserSegmentRule
}

func NewUserSegmentRuleService() realtimecore.SegmentRuleService {
	return &userSegmentRuleService{
		rules: map[string]realtimecore.UserSegmentRule{},
	}
}

func (s *userSegmentRuleService) Put(ctx context.Context, SegmentRule realtimecore.UserSegmentRule) error {
	s.Lock()
	defer s.Unlock()

	s.rules[SegmentRule.Name] = SegmentRule
	return nil
}

func (s *userSegmentRuleService) Get(ctx context.Context, segmentUID string) (realtimecore.UserSegmentRule, bool, error) {
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

func (s *userSegmentRuleService) List(ctx context.Context) ([]realtimecore.UserSegmentRule, error) {
	s.Lock()
	defer s.Unlock()

	segments := []realtimecore.UserSegmentRule{}

	for _, sd := range s.rules {
		segments = append(segments, sd)
	}

	return segments, nil
}
