package datastore

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoTransaction(t *testing.T) {
	ctx := context.TODO()
	i := 42

	t.Run("get, empty", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given

		// when
		_, exists, err := store.Get(ctx, "1")

		// then
		assert.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("put", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given

		// when
		err := store.Put(ctx, "1", 1)

		// then
		assert.NoError(t, err)
		val, exists, err := store.Get(ctx, "1")
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, 1, val)
	})

	t.Run("put wrong data-type", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given

		// when
		err := store.Put(ctx, "1", "string")

		// then
		assert.Error(t, err)
		assert.Equal(t, "Unexpected data-type '.string': want '.int'", err.Error())

	})

	t.Run("remove", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given
		store.Put(ctx, "1", 1)

		// when
		err := store.Remove(ctx, "1")

		// then
		assert.NoError(t, err)
		_, exists, err := store.Get(ctx, "1")
		assert.False(t, exists)

	})

	t.Run("get all", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given
		store.Put(ctx, "1", 1)

		// when
		result, err := store.GetAll(ctx)

		// then
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 1, result[0])
	})

}

func TestTransaction(t *testing.T) {
	ctx := context.TODO()
	i := 42

	t.Run("get, empty", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given

		// when
		err := store.RunInTransaction(ctx, func(ctx context.Context) error {
			_, exists, err := store.Get(ctx, "1")

			// then
			assert.NoError(t, err)
			assert.False(t, exists)
			return nil
		})

		// then
		assert.NoError(t, err)
	})

	t.Run("put", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given

		// when
		err := store.RunInTransaction(ctx, func(ctx context.Context) error {
			err := store.Put(ctx, "1", 1)
			// Not yet visible in transaction
			_, exists, _ := store.get(ctx, "1")
			assert.False(t, exists)
			return err
		})

		// then
		assert.NoError(t, err)
		val, exists, err := store.Get(ctx, "1")
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.Equal(t, 1, val)
	})

	t.Run("remove", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given
		store.Put(ctx, "1", 1)

		// when
		err := store.RunInTransaction(ctx, func(ctx context.Context) error {
			err := store.Remove(ctx, "1")
			// Still visible in transaction
			_, exists, _ := store.get(ctx, "1")
			assert.True(t, exists)
			return err
		})

		// then
		assert.NoError(t, err)
		_, exists, _ := store.Get(ctx, "1")
		assert.False(t, exists)
	})

	t.Run("get all", func(t *testing.T) {
		store := NewDatastoreStub()
		store.EnforceDataType(i)

		// given
		store.Put(ctx, "1", 1)

		// when
		var result []interface{}
		var err error
		err = store.RunInTransaction(ctx, func(ctx context.Context) error {
			result, err = store.GetAll(ctx)
			return err
		})

		// then
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, 1, result[0])
	})

}
