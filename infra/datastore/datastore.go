package datastore

import (
	"context"
	"sync"
)

type lockedInMemoryDatastore struct {
	sync.Mutex // Poor-mans transaction
	items      map[string]interface{}
}

func NewDatastore() Datastore {
	return &lockedInMemoryDatastore{
		items: map[string]interface{}{},
	}
}

func (s *lockedInMemoryDatastore) RunInTransaction(ctx context.Context, callback func(ctx context.Context) error) error {
	// start transaction
	s.Lock()

	err := callback(ctx)
	if err != nil {
		// abort transaction
		s.Unlock()
		return err
	}

	// commit transaction
	s.Unlock()

	return nil
}

func (s *lockedInMemoryDatastore) Put(ctx context.Context, uid string, item interface{}) error {
	s.items[uid] = item

	return nil
}

func (s *lockedInMemoryDatastore) Remove(ctx context.Context, uid string) error {
	_, exists := s.items[uid]
	if exists {
		delete(s.items, uid)
	}

	return nil
}

func (s *lockedInMemoryDatastore) Get(ctx context.Context, uid string) (interface{}, bool, error) {
	user, exists := s.items[uid]
	return user, exists, nil
}

func (s *lockedInMemoryDatastore) GetAll(ctx context.Context) ([]interface{}, error) {
	result := []interface{}{}
	for _, i := range s.items {
		result = append(result, i)
	}

	return result, nil
}
