package segment

import (
	"context"
)

type SegmentManagementStub struct {
	Segments map[string]SegmentSpec
}

func NewSegmentManagementStub() *SegmentManagementStub {
	return &SegmentManagementStub{
		Segments: map[string]SegmentSpec{},
	}
}

func (s *SegmentManagementStub) Put(ctx context.Context, u SegmentSpec) error {
	s.Segments[u.UID] = u
	return nil
}

func (s *SegmentManagementStub) Remove(ctx context.Context, uid string) error {
	delete(s.Segments, uid)
	return nil
}

func (s *SegmentManagementStub) Get(ctx context.Context, uid string) (SegmentSpec, bool, error) {
	item, exists := s.Segments[uid]
	return item, exists, nil
}
func (s *SegmentManagementStub) List(ctx context.Context) ([]SegmentSpec, error) {
	items := []SegmentSpec{}
	for _, i := range s.Segments {
		items = append(items, i)
	}
	return items, nil
}
