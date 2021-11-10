package usermanagement

import (
	"context"
	"github.com/MarcGrol/userautomation/coredata/supportedattrs"
	"github.com/MarcGrol/userautomation/services/filtermanagement"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/coredata/predefinedusers"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestUserManagement(t *testing.T) {
	ctx := context.TODO()

	t.Run("create user", func(t *testing.T) {
		// setup

		userStore, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(userStore, filterService, ps)

		// given
		nothing()

		// when
		err := sut.Put(ctx, initialUser())

		// then
		assert.NoError(t, err)
		assert.Len(t, listUser(ctx, t, sut), 1)
		assert.Equal(t, "Marc", getUser(ctx, t, sut).Attributes[supportedattrs.FirstName])
		assert.Equal(t, initialUser().UID, ps.Publications[0].Event.(user.CreatedEvent).UserState.UID)
	})

	t.Run("modify user", func(t *testing.T) {
		// setup
		userStore, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(userStore, filterService, ps)

		// given
		sut.Put(ctx, initialUser())

		// when
		err := sut.Put(ctx, modifiedUser())

		// then
		assert.NoError(t, err)
		assert.Len(t, listUser(ctx, t, sut), 1)
		assert.Equal(t, "Eva", getUser(ctx, t, sut).Attributes[supportedattrs.FirstName])
		assert.Equal(t, "Marc", ps.Publications[1].Event.(user.ModifiedEvent).OldUserState.Attributes[supportedattrs.FirstName])
		assert.Equal(t, "Eva", ps.Publications[1].Event.(user.ModifiedEvent).NewUserState.Attributes[supportedattrs.FirstName])
	})

	t.Run("remove user", func(t *testing.T) {
		// setup
		userStore, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(userStore, filterService, ps)

		// given
		sut.Put(ctx, initialUser())

		// when
		err := sut.Remove(ctx, initialUser().UID)

		// then
		assert.NoError(t, err)
		assert.False(t, existUser(ctx, t, sut))
		assert.Len(t, listUser(ctx, t, sut), 0)
		assert.Equal(t, initialUser().UID, ps.Publications[1].Event.(user.RemovedEvent).UserState.UID)
	})
}

func setupMocks(t *testing.T) (*datastore.DatastoreStub, user.FilterManager, *pubsub.PubsubStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	storeSub := datastore.NewDatastoreStub()
	ps := pubsub.NewPubsubStub()
	filterManager := filtermanagement.New()
	return storeSub, filterManager, ps, ctrl
}

func initialUser() user.User {
	return predefinedusers.Marc
}

func modifiedUser() user.User {
	u := predefinedusers.Eva
	u.UID = predefinedusers.Marc.UID
	return u
}

func existUser(ctx context.Context, t *testing.T, sut user.Management) bool {
	_, exists, err := sut.Get(ctx, initialUser().UID)
	if err != nil {
		t.Error(err)
	}
	return exists
}

func getUser(ctx context.Context, t *testing.T, sut user.Management) user.User {
	u, exists, err := sut.Get(ctx, initialUser().UID)
	if err != nil || !exists {
		t.Error(err)
	}
	return u
}

func listUser(ctx context.Context, t *testing.T, sut user.Management) []user.User {
	users, err := sut.List(ctx)
	if err != nil {
		t.Error(err)
	}
	return users
}

func nothing() {}
