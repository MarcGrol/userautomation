package datastore

import (
	"context"
	"sync"
)

type lockedInMemDatastore struct {
	sync.Mutex
	items map[string]interface{}
}

func NewDatastore() Datastore {
	return &lockedInMemDatastore{
		items: map[string]interface{}{},
	}
}

func (s *lockedInMemDatastore) RunInTransaction(ctx context.Context, callback func(ctx context.Context) error) error {
	s.Lock()
	defer s.Unlock()

	err := callback(ctx)
	if err != nil {
		// abort transaction
		return err
	}
	// commit transaction

	return nil
}

func (s *lockedInMemDatastore) Put(ctx context.Context, uid string, item interface{}) error {
	s.items[uid] = item

	return nil
}

func (s *lockedInMemDatastore) Remove(ctx context.Context, uid string) error {
	_, exists := s.items[uid]
	if exists {
		// remove from memory
		delete(s.items, uid)
	}

	return nil
}

func (s *lockedInMemDatastore) Get(ctx context.Context, uid string) (interface{}, bool, error) {
	user, exists := s.items[uid]
	return user, exists, nil
}

func (s *lockedInMemDatastore) GetAll(ctx context.Context) ([]interface{}, error) {
	result := []interface{}{}
	for _, i := range s.items {
		result = append(result, i)
	}

	return result, nil
}
