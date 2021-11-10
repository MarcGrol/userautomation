package segment

import (
	"context"
)

type SegmentManagementStub struct {
	Segments map[string]Spec
}

func NewSegmentManagementStub() *SegmentManagementStub {
	return &SegmentManagementStub{
		Segments: map[string]Spec{},
	}
}

func (s *SegmentManagementStub) Put(ctx context.Context, u Spec) error {
	s.Segments[u.UID] = u
	return nil
}

func (s *SegmentManagementStub) Remove(ctx context.Context, uid string) error {
	delete(s.Segments, uid)
	return nil
}

func (s *SegmentManagementStub) Get(ctx context.Context, uid string) (Spec, bool, error) {
	item, exists := s.Segments[uid]
	return item, exists, nil
}
func (s *SegmentManagementStub) List(ctx context.Context) ([]Spec, error) {
	items := []Spec{}
	for _, i := range s.Segments {
		items = append(items, i)
	}
	return items, nil
}
