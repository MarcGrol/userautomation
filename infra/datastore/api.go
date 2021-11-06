package datastore

import (
	"context"
)

//go:generate mockgen -source=api.go -destination=datastore_mock.go -package=datastore Datastore
type Datastore interface {
	RunInTransaction(ctx context.Context, callback func(ctx context.Context) error) error
	Put(ctx context.Context, uid string, item interface{}) error
	Remove(ctx context.Context, uid string) error
	Get(ctx context.Context, uid string) (interface{}, bool, error)
	GetAll(ctx context.Context) ([]interface{}, error)
}
