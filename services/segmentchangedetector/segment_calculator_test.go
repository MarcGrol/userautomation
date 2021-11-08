package segmentchangedetector

import (
	"context"
	"log"
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/MarcGrol/userautomation/core/segment"
	"github.com/MarcGrol/userautomation/core/user"
	"github.com/MarcGrol/userautomation/infra/datastore"
	"github.com/MarcGrol/userautomation/infra/pubsub"
	"github.com/stretchr/testify/assert"
)

func TestSegment(t *testing.T) {
	ctx := context.TODO()

	t.Run("on segment created, no users", func(t *testing.T) {
		// setup

		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		noUsers()

		// when
		onSegmentCreated(ctx, t, sut)

		// then
		assert.Equal(t, "young", getSegment(ctx, t, sut).Description)
		assert.Empty(t, getSegment(ctx, t, sut).Users)
	})

	t.Run("on segment created, with non matching users", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 50)

		// when
		onSegmentCreated(ctx, t, sut)

		// then
		assert.Equal(t, "young", getSegment(ctx, t, sut).Description)
		assert.Empty(t, getSegment(ctx, t, sut).Users)
	})

	t.Run("on segment created, with matching user", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 12)

		// when
		onSegmentCreated(ctx, t, sut)

		// then
		assert.Equal(t, "young", getSegment(ctx, t, sut).Description)
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)
	})

	t.Run("on segment modified, segment with two old users", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 50)
		createOtherUser(ctx, t, userService, 47)
		createEmptySegment(ctx, t, sut)

		// when
		onSegmentModified(ctx, t, sut)

		// then
		log.Printf("segment%+v", getSegment(ctx, t, sut))
		assert.Equal(t, "old", getSegment(ctx, t, sut).Description)
		assert.Len(t, getSegment(ctx, t, sut).Users, 2)
	})

	t.Run("on-segment-removed", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 12)
		createOtherUser(ctx, t, userService, 12)
		createEmptySegment(ctx, t, sut)

		// when
		err := onSegmentRemoved(ctx, t, sut)

		// then
		assert.NoError(t, err)
	})

	t.Run("user-created, no segment exists", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		noUsers()

		// when
		sut.OnUserCreated(ctx, user.CreatedEvent{UserState: getUser(50)})

		// then
	})

	t.Run("user-created, no matching segment exists", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		noUsers()
		onSegmentCreated(ctx, t, sut)

		// when
		sut.OnUserCreated(ctx, user.CreatedEvent{UserState: getUser(50)})

		// then
	})

	t.Run("user-created, matching segment: user added to segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		noUsers()
		createEmptySegment(ctx, t, sut)

		// when
		sut.OnUserCreated(ctx, user.CreatedEvent{UserState: getUser(12)})

		// then
		assert.True(t, existsUserInSegment(ctx, t, sut, "1"))
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)
	})

	t.Run("user-modified, matching segment: user added to segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given

		createUser(ctx, t, userService, 50)
		createEmptySegment(ctx, t, sut)
		segm1, _, _ := sut.segmentWithUsersStore.Get(ctx, "YoungSegment")
		log.Printf("segm:%+v", segm1)

		segm, _, _ := sut.segmentWithUsersStore.Get(ctx, "YoungSegment")
		log.Printf("segm:%+v", segm)
		assert.Len(t, getSegment(ctx, t, sut).Users, 0)

		// when
		sut.OnUserModified(ctx, user.ModifiedEvent{OldUserState: getUser(50), NewUserState: getUser(12)})

		// then
		assert.True(t, existsUserInSegment(ctx, t, sut, "1"))
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)
	})

	t.Run("user-modified, no matching segment: user removed from segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		u := createUser(ctx, t, userService, 12)
		createSegmentWithUsers(ctx, t, sut, u)
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)

		// when
		sut.OnUserModified(ctx, user.ModifiedEvent{OldUserState: getUser(12), NewUserState: getUser(50)})

		// then
		assert.Len(t, getSegment(ctx, t, sut).Users, 0)
	})

	t.Run("user-removed, matching segment: user removed from segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		u := createUser(ctx, t, userService, 13)
		createSegmentWithUsers(ctx, t, sut, u)
		assert.Len(t, getSegment(ctx, t, sut).Users, 1)

		// when
		sut.OnUserRemoved(ctx, user.RemovedEvent{UserState: getUser(12)})

		// then
		assert.Len(t, getSegment(ctx, t, sut).Users, 0)
	})

}

func setupMocks(t *testing.T) (*datastore.DatastoreStub, *user.UserManagementStub, *pubsub.PubsubStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	segmentStore := datastore.NewDatastoreStub()
	userService := user.NewUserManagementStub()
	ps := pubsub.NewPubsubStub()
	return segmentStore, userService, ps, ctrl
}

var initialSegment = func() segment.UserSegment {
	return segment.UserSegment{
		UID:            "YoungSegment",
		Description:    "young",
		UserFilterName: user.FilterYoungAge,
		Users:          map[string]user.User{},
	}
}

func createEmptySegment(ctx context.Context, t *testing.T, sut *segmentCalculator) {
	segm := initialSegment()
	err := sut.segmentWithUsersStore.Put(ctx, segm.UID, segm)
	if err != nil {
		t.Error(err)
	}
}

func createSegmentWithUsers(ctx context.Context, t *testing.T, sut *segmentCalculator, users ...user.User) {
	segm := initialSegment()

	for _, u := range users {
		segm.Users[u.UID] = u
	}
	err := sut.segmentWithUsersStore.Put(ctx, segm.UID, segm)
	if err != nil {
		t.Error(err)
	}
}

func onSegmentCreated(ctx context.Context, t *testing.T, sut segment.EventHandler) {
	event := segment.CreatedEvent{
		SegmentState: initialSegment(),
	}
	err := sut.OnSegmentCreated(ctx, event)
	if err != nil {
		t.Error(err)
	}
}

var modifiedSegment = func() segment.UserSegment {
	return segment.UserSegment{
		UID:            "YoungSegment",
		Description:    "old",
		UserFilterName: user.FilterOldAge,
		Users:          map[string]user.User{},
	}
}

func onSegmentModified(ctx context.Context, t *testing.T, sut segment.EventHandler) error {
	event := segment.ModifiedEvent{
		OldSegmentState: initialSegment(),
		NewSegmentState: modifiedSegment(),
	}
	return sut.OnSegmentModified(ctx, event)
}

func onSegmentRemoved(ctx context.Context, t *testing.T, sut segment.EventHandler) error {
	event := segment.RemovedEvent{
		SegmentState: initialSegment(),
	}
	return sut.OnSegmentRemoved(ctx, event)
}

func getSegment(ctx context.Context, t *testing.T, sut *segmentCalculator) segment.UserSegment {
	segm, exists, err := sut.segmentWithUsersStore.Get(ctx, "YoungSegment")
	if err != nil || !exists {
		t.Error(err)
	}
	return segm.(segment.UserSegment)
}

func existsUserInSegment(ctx context.Context, t *testing.T, sut *segmentCalculator, userId string) bool {
	segm := getSegment(ctx, t, sut)
	_, exists := segm.Users[userId]
	return exists
}

func noUsers() {}

func createUser(ctx context.Context, t *testing.T, userService user.Management, age int) user.User {
	u := user.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          age,
		},
	}
	err := userService.Put(ctx, u)
	if err != nil {
		t.Error(err)
	}
	return u
}

func createOtherUser(ctx context.Context, t *testing.T, userService user.Management, age int) {
	err := userService.Put(ctx, user.User{
		UID: "2",
		Attributes: map[string]interface{}{
			"firstname":    "Eva",
			"emailaddress": "eva@home.nl",
			"phonenumber":  "+31622222222",
			"age":          age,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func getUser(age int) user.User {
	return user.User{
		UID: "1",
		Attributes: map[string]interface{}{
			"firstname":    "Marc",
			"emailaddress": "marc@home.nl",
			"phonenumber":  "+31611111111",
			"age":          age,
		},
	}
}

func getOtherUser(age int) user.User {
	return user.User{
		UID: "2",
		Attributes: map[string]interface{}{
			"firstname":    "Eva",
			"emailaddress": "eva@home.nl",
			"phonenumber":  "+31622222222",
			"age":          age,
		},
	}
}
