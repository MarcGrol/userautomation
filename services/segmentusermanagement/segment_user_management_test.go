package segmentusermanagement

import (
	"context"
	"github.com/MarcGrol/userautomation/coredata/predefinedfilters"
	"github.com/MarcGrol/userautomation/services/filtermanager"
	"testing"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/coredata/predefinedsegments"
	"github.com/MarcGrol/userautomation/coredata/predefinedusers"
	"github.com/MarcGrol/userautomation/coredata/supportedattrs"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestSegmentUserManagement(t *testing.T) {
	ctx := context.TODO()

	t.Run("on segment created, no users", func(t *testing.T) {
		// setup

		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		noUsers()

		// when
		err := sut.OnSegmentCreated(ctx, segment.CreatedEvent{
			SegmentState: initialSegmentWithUsers().SegmentSpec,
		})

		// then
		assert.NoError(t, err)
		assert.Equal(t, "young users segment", getSegment(ctx, t, sut).SegmentSpec.Description)
		assert.Empty(t, getSegment(ctx, t, sut).Users)
		assert.Len(t, ps.Publications, 0)
	})

	t.Run("on segment created, with non matching users", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		createUserWithAge(ctx, t, userService, 50)

		// when
		err := sut.OnSegmentCreated(ctx, segment.CreatedEvent{
			SegmentState: initialSegmentWithUsers().SegmentSpec,
		})

		// then
		assert.NoError(t, err)
		assert.Equal(t, "young users segment", getSegment(ctx, t, sut).SegmentSpec.Description)
		assert.Empty(t, getSegment(ctx, t, sut).Users)
		assert.Len(t, ps.Publications, 0)
	})

	t.Run("on segment created, with matching user", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		createUserWithAge(ctx, t, userService, 12)

		// when
		err := sut.OnSegmentCreated(ctx, segment.CreatedEvent{
			SegmentState: initialSegmentWithUsers().SegmentSpec,
		})

		// then
		assert.NoError(t, err)
		assert.Equal(t, "young users segment", getSegment(ctx, t, sut).SegmentSpec.Description)
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)
		assert.Len(t, ps.Publications, 1)
		assert.Equal(t, defaultUser().UID, ps.Publications[0].Event.(segment.UserAddedToSegmentEvent).User.UID)
	})

	t.Run("on segment modified, segment with two old users", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		createUserWithAge(ctx, t, userService, 50)
		createOtherUser(ctx, t, userService, 47)
		createEmptySegmentWithUsers(ctx, t, sut)

		// when
		err := sut.OnSegmentModified(ctx, segment.ModifiedEvent{
			OldSegmentState: initialSegmentWithUsers().SegmentSpec,
			NewSegmentState: modifiedSegment().SegmentSpec,
		})

		// then
		assert.NoError(t, err)
		assert.Equal(t, "old", getSegment(ctx, t, sut).SegmentSpec.Description)
		assert.Len(t, getSegment(ctx, t, sut).Users, 2)
		assert.Len(t, ps.Publications, 2)
		assert.Equal(t, defaultUser().UID, ps.Publications[0].Event.(segment.UserAddedToSegmentEvent).User.UID)
		assert.Equal(t, otherUser().UID, ps.Publications[1].Event.(segment.UserAddedToSegmentEvent).User.UID)
	})

	t.Run("on-segment-removed", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		createUserWithAge(ctx, t, userService, 12)
		createOtherUser(ctx, t, userService, 12)
		createEmptySegmentWithUsers(ctx, t, sut)

		// when
		err := sut.OnSegmentRemoved(ctx, segment.RemovedEvent{
			SegmentState: initialSegmentWithUsers().SegmentSpec,
		})

		// then
		assert.NoError(t, err)
		assert.Len(t, ps.Publications, 0)
	})

	t.Run("user-created, no segment exists", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		noUsers()

		// when
		err := sut.OnUserCreated(ctx, user.CreatedEvent{UserState: userWithAge(50)})

		// then
		assert.NoError(t, err)
		assert.Len(t, ps.Publications, 0)
	})

	t.Run("user-created, no matching segment exists", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		noUsers()
		sut.OnSegmentCreated(ctx, segment.CreatedEvent{
			SegmentState: initialSegmentWithUsers().SegmentSpec,
		})

		// when
		err := sut.OnUserCreated(ctx, user.CreatedEvent{UserState: userWithAge(50)})

		// then
		assert.NoError(t, err)
		assert.Len(t, ps.Publications, 0)
	})

	t.Run("user-created, matching segment: user added to segment", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		noUsers()
		createEmptySegmentWithUsers(ctx, t, sut)

		// when
		err := sut.OnUserCreated(ctx, user.CreatedEvent{UserState: userWithAge(12)})

		// then
		assert.NoError(t, err)
		assert.True(t, existsUserInSegment(ctx, t, sut, "1"))
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)
		assert.Len(t, ps.Publications, 1)
		assert.Equal(t, defaultUser().UID, ps.Publications[0].Event.(segment.UserAddedToSegmentEvent).User.UID)
	})

	t.Run("user-modified, matching segment: user added to segment", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		createUserWithAge(ctx, t, userService, 50)
		createEmptySegmentWithUsers(ctx, t, sut)
		assert.Len(t, getSegment(ctx, t, sut).Users, 0)

		// when
		err := sut.OnUserModified(ctx, user.ModifiedEvent{OldUserState: userWithAge(50), NewUserState: userWithAge(12)})

		// then
		assert.NoError(t, err)
		assert.True(t, existsUserInSegment(ctx, t, sut, "1"))
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)
		assert.Len(t, ps.Publications, 1)
		assert.Equal(t, defaultUser().UID, ps.Publications[0].Event.(segment.UserAddedToSegmentEvent).User.UID)
	})

	t.Run("user-modified, matching segment: user removed from segment", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		u := createUserWithAge(ctx, t, userService, 12)
		createSegmentWithUsers(ctx, t, sut, u)
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)

		// when
		err := sut.OnUserModified(ctx, user.ModifiedEvent{OldUserState: userWithAge(12), NewUserState: userWithAge(50)})

		// then
		assert.NoError(t, err)
		assert.Len(t, getSegment(ctx, t, sut).Users, 0)
		assert.Len(t, ps.Publications, 1)
		assert.Equal(t, defaultUser().UID, ps.Publications[0].Event.(segment.UserRemovedFromSegmentEvent).User.UID)
	})

	t.Run("user-removed, matching segment: user removed from segment", func(t *testing.T) {
		// setup
		segmentStore, userService, filterService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, filterService, ps)

		// given
		u := createUserWithAge(ctx, t, userService, 13)
		createSegmentWithUsers(ctx, t, sut, u)
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)

		// when
		err := sut.OnUserRemoved(ctx, user.RemovedEvent{UserState: userWithAge(12)})

		// then
		assert.NoError(t, err)
		assert.Len(t, getSegment(ctx, t, sut).Users, 0)
		assert.Len(t, ps.Publications, 1)
		assert.Equal(t, defaultUser().UID, ps.Publications[0].Event.(segment.UserRemovedFromSegmentEvent).User.UID)
	})

}

func setupMocks(t *testing.T) (*datastore.DatastoreStub, *user.UserManagementStub, user.FilterManager, *pubsub.PubsubStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	segmentStore := datastore.NewDatastoreStub()
	filterService := filtermanager.New()
	userService := user.NewUserManagementStub(filterService)
	ps := pubsub.NewPubsubStub()
	return segmentStore, userService, filterService, ps, ctrl
}

func initialSegmentWithUsers() segment.WithUsers {
	return segment.WithUsers{
		SegmentSpec: predefinedsegments.YoungAgeSegment,
		Users:       map[string]user.User{},
	}
}

func createEmptySegmentWithUsers(ctx context.Context, t *testing.T, sut *segmentUserManager) {
	swu := initialSegmentWithUsers()
	err := sut.segmentWithUsersStore.Put(ctx, swu.SegmentSpec.UID, swu)
	if err != nil {
		t.Error(err)
	}
}

func createSegmentWithUsers(ctx context.Context, t *testing.T, sut *segmentUserManager, users ...user.User) {
	swu := initialSegmentWithUsers()

	for _, u := range users {
		swu.Users[u.UID] = u
	}
	err := sut.segmentWithUsersStore.Put(ctx, swu.SegmentSpec.UID, swu)
	if err != nil {
		t.Error(err)
	}
}

func modifiedSegment() segment.WithUsers {
	swu := initialSegmentWithUsers()
	swu.SegmentSpec.Description = "old"
	swu.SegmentSpec.UserFilterName = predefinedfilters.FilterOldAgeName
	swu.Users = map[string]user.User{}

	return swu
}

func getSegment(ctx context.Context, t *testing.T, sut *segmentUserManager) segment.WithUsers {
	item, exists, err := sut.segmentWithUsersStore.Get(ctx, predefinedsegments.YoungAgeSegmentName)
	if err != nil || !exists {
		t.Error(err)
	}
	return item.(segment.WithUsers)
}

func existsUserInSegment(ctx context.Context, t *testing.T, sut *segmentUserManager, userId string) bool {
	swu := getSegment(ctx, t, sut)
	_, exists := swu.Users[userId]
	return exists
}

func noUsers() {}

func defaultUser() user.User {
	return predefinedusers.Marc
}

func createUserWithAge(ctx context.Context, t *testing.T, userService user.Management, age int) user.User {
	u := defaultUser()
	u.Attributes[supportedattrs.Age] = age
	err := userService.Put(ctx, u)
	if err != nil {
		t.Error(err)
	}
	return u
}

func userWithAge(age int) user.User {
	u := defaultUser()
	u.Attributes[supportedattrs.Age] = age
	return u
}

func otherUser() user.User {
	return predefinedusers.Eva
}

func createOtherUser(ctx context.Context, t *testing.T, userService user.Management, age int) {
	err := userService.Put(ctx, otherUser())
	if err != nil {
		t.Error(err)
	}
}
