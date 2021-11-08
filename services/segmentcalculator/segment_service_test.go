package segmentcalculator

import (
	"context"
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

	t.Run("create segment, no users", func(t *testing.T) {
		// setup

		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		noUsers()

		// when
		createYoungSegment(ctx, t, sut, "x")

		// then
		assert.Equal(t, "x", getYoungSegment(ctx, t, sut).Description)
		assert.Empty(t, getYoungSegment(ctx, t, sut).Users)
	})

	t.Run("create segment, with non matching users", func(t *testing.T) {
		// setup

		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 50)

		// when
		createYoungSegment(ctx, t, sut, "x")

		// then
		assert.Equal(t, "x", getYoungSegment(ctx, t, sut).Description)
		assert.Empty(t, getYoungSegment(ctx, t, sut).Users)
	})

	t.Run("create segment, with matching user", func(t *testing.T) {
		// setup

		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 12)

		// when
		createYoungSegment(ctx, t, sut, "x")

		// then
		assert.Equal(t, "x", getYoungSegment(ctx, t, sut).Description)
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 1)
	})

	t.Run("modify segment with two young users", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 12)
		createOtherUser(ctx, t, userService, 12)
		createYoungSegment(ctx, t, sut, "x")

		// when
		createYoungSegment(ctx, t, sut, "y")

		// then
		assert.Equal(t, "y", getYoungSegment(ctx, t, sut).Description)
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 2)
	})

	t.Run("remove segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 12)
		createOtherUser(ctx, t, userService, 12)
		createYoungSegment(ctx, t, sut, "x")

		// when
		err := removeYoungSegment(ctx, t, sut)

		// then
		assert.Error(t, err)
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
		createYoungSegment(ctx, t, sut, "x")

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
		createYoungSegment(ctx, t, sut, "x")
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 0)

		// when
		sut.OnUserCreated(ctx, user.CreatedEvent{UserState: getUser(12)})

		// then
		assert.True(t, existsUserInYoungSegment(ctx, t, sut, "1"))
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 1)
	})

	t.Run("user-modified, matching segment: user added to segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 50)
		createYoungSegment(ctx, t, sut, "x")
		sut.OnUserCreated(ctx, user.CreatedEvent{getUser(50)})
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 0)

		// when
		sut.OnUserModified(ctx, user.ModifiedEvent{OldUserState: getUser(50), NewUserState: getUser(12)})

		// then
		assert.True(t, existsUserInYoungSegment(ctx, t, sut, "1"))
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 1)
	})

	t.Run("user-modified, no matching segment: user removed from segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 12)
		createYoungSegment(ctx, t, sut, "x")
		sut.OnUserCreated(ctx, user.CreatedEvent{getUser(12)})
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 1)

		// when
		sut.OnUserModified(ctx, user.ModifiedEvent{OldUserState: getUser(12), NewUserState: getUser(50)})

		// then
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 0)
	})

	t.Run("user-removed, matching segment: user removed from segment", func(t *testing.T) {
		// setup
		segmentStore, userService, ps, ctrl := setupMocks(t)
		defer ctrl.Finish()
		sut := New(segmentStore, userService, ps)

		// given
		createUser(ctx, t, userService, 13)
		createYoungSegment(ctx, t, sut, "x")
		sut.OnUserCreated(ctx, user.CreatedEvent{UserState: getUser(12)})
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 1)

		// when
		sut.OnUserRemoved(ctx, user.RemovedEvent{UserState: getUser(12)})

		// then
		assert.Len(t, getYoungSegment(ctx, t, sut).Users, 0)
	})

}

func setupMocks(t *testing.T) (*datastore.DatastoreStub, *user.UserServiceStub, *pubsub.PubsubStub, *gomock.Controller) {
	ctrl := gomock.NewController(t)
	segmentStore := datastore.NewDatastoreStub()
	userService := user.NewUserServiceStub()
	ps := pubsub.NewPubsubStub()
	return segmentStore, userService, ps, ctrl
}

func createYoungSegment(ctx context.Context, t *testing.T, sut SegmentService, description string) {
	segment := segment.UserSegment{
		UID:         "YoungSegment",
		Description: description,
		IsApplicableForUser: func(ctx context.Context, user user.User) (bool, error) {
			age, ok := user.Attributes["age"].(int)
			if !ok {
				return false, nil
			}
			return age < 18, nil
		},
	}
	err := sut.Put(ctx, segment)
	if err != nil {
		t.Error(err)
	}
}

func removeYoungSegment(ctx context.Context, t *testing.T, sut SegmentService) error {
	err := sut.Remove(ctx, "YoungSegment")
	if err != nil {
		return err
	}
	return nil
}

func getYoungSegment(ctx context.Context, t *testing.T, sut SegmentService) segment.UserSegment {
	segm, exists, err := sut.Get(ctx, "YoungSegment")
	if err != nil || !exists {
		t.Error(err)
	}
	return segm
}

func existsUserInYoungSegment(ctx context.Context, t *testing.T, sut SegmentService, userId string) bool {
	segm := getYoungSegment(ctx, t, sut)
	_, exists := segm.Users[userId]
	return exists
}

func noUsers() {}

func createUser(ctx context.Context, t *testing.T, userService user.Service, age int) user.User {
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

func createOtherUser(ctx context.Context, t *testing.T, userService user.Service, age int) {
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
