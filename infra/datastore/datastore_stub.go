package datastore

import (
	"context"
	"sync"
)

// TODO integrate with 3rd party persistent datastore

const (
	transactionContextKey = "transaction"
)

type DatastoreStub struct {
	sync.Mutex
	itemKind string
	Items    map[string]interface{}
}

func NewDatastoreStub() *DatastoreStub {
	return &DatastoreStub{
		Items: map[string]interface{}{},
	}
}

func (s *DatastoreStub) RunInTransaction(ctx context.Context, callback func(ctx context.Context) error) error {
	s.Lock()
	defer s.Unlock()

	// Start transaction
	trx := &transaction{
		operations: []datastoreOperation{},
	}

	// Call inner logic with enriched, transaction-aware context
	err := callback(context.WithValue(ctx, transactionContextKey, trx))
	if err != nil {
		// Abort transaction
		return err
	}

	// Commit transaction
	err = s.commitPendingOperations(ctx, trx)
	if err != nil {
		return err
	}

	return nil
}

func (s *DatastoreStub) commitPendingOperations(ctx context.Context, trx *transaction) error {
	for _, oper := range trx.operations {
		switch oper.Kind {
		case datastoreOperationUpsert:
			err := s.put(ctx, oper.UID, oper.data)
			if err != nil {
				return err
			}
		case datastoreOperationRemove:
			err := s.remove(ctx, oper.UID)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *DatastoreStub) Put(ctx context.Context, uid string, item interface{}) error {
	t := ctx.Value(transactionContextKey)
	if t != nil {
		trx := t.(*transaction)
		return trx.put(uid, item)
	}

	s.Lock()
	defer s.Unlock()
	return s.put(ctx, uid, item)
}

func (s *DatastoreStub) put(_ context.Context, uid string, item interface{}) error {
	s.Items[uid] = item

	return nil
}

func (s *DatastoreStub) Remove(ctx context.Context, uid string) error {
	t := ctx.Value(transactionContextKey)
	if t != nil {
		trx := t.(*transaction)
		return trx.remove(uid)
	}

	s.Lock()
	defer s.Unlock()
	return s.remove(ctx, uid)
}

func (s *DatastoreStub) remove(_ context.Context, uid string) error {
	_, exists := s.Items[uid]
	if exists {
		delete(s.Items, uid)
	}

	return nil
}

func (s *DatastoreStub) Get(ctx context.Context, uid string) (interface{}, bool, error) {
	t := ctx.Value(transactionContextKey)
	if t != nil {
		// Already locked by RunInTransaction
		return s.get(ctx, uid)
	}

	s.Lock()
	defer s.Unlock()
	return s.get(ctx, uid)
}

func (s *DatastoreStub) get(ctx context.Context, uid string) (interface{}, bool, error) {
	user, exists := s.Items[uid]
	return user, exists, nil
}

func (s *DatastoreStub) GetAll(ctx context.Context) ([]interface{}, error) {
	t := ctx.Value(transactionContextKey)
	if t != nil {
		// Already locked by RunInTransaction
		return s.getAll(ctx)
	}

	s.Lock()
	defer s.Unlock()
	return s.getAll(ctx)
}

func (s *DatastoreStub) getAll(ctx context.Context) ([]interface{}, error) {
	result := []interface{}{}
	for _, i := range s.Items {
		result = append(result, i)
	}

	return result, nil
}
