package segment

import (
	"context"
)

type SegmentManagementStub struct {
	Segments map[string]UserSegment
}

func NewSegmentManagementStub() *SegmentManagementStub {
	return &SegmentManagementStub{
		Segments: map[string]UserSegment{},
	}
}

func (s *SegmentManagementStub) Put(ctx context.Context, u UserSegment) error {
	s.Segments[u.UID] = u
	return nil
}

func (s *SegmentManagementStub) Remove(ctx context.Context, uid string) error {
	delete(s.Segments, uid)
	return nil
}

func (s *SegmentManagementStub) Get(ctx context.Context, uid string) (UserSegment, bool, error) {
	item, exists := s.Segments[uid]
	return item, exists, nil
}
func (s *SegmentManagementStub) List(ctx context.Context) ([]UserSegment, error) {
	items := []UserSegment{}
	for _, i := range s.Segments {
		items = append(items, i)
	}
	return items, nil
}
