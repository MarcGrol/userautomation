package datastore

import (
	"context"
)

type DatastoreStub struct {
	Items map[string]interface{}
}

func NewDatastoreStub() *DatastoreStub {
	return &DatastoreStub{
		Items: map[string]interface{}{},
	}
}

func (s *DatastoreStub) RunInTransaction(ctx context.Context, callback func(ctx context.Context) error) error {
	return callback(ctx)
}

func (s *DatastoreStub) Put(ctx context.Context, uid string, item interface{}) error {
	s.Items[uid] = item
	return nil
}

func (s *DatastoreStub) Remove(ctx context.Context, uid string) error {
	delete(s.Items, uid)
	return nil
}

func (s *DatastoreStub) Get(ctx context.Context, uid string) (interface{}, bool, error) {
	item, exists := s.Items[uid]
	return item, exists, nil
}

func (s *DatastoreStub) GetAll(ctx context.Context) ([]interface{}, error) {
	items := []interface{}{}
	for _, i := range s.Items {
		items = append(items, i)
	}
	return items, nil
}
